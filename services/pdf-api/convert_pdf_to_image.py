# create a kafka consumer
import json

from kafka import KafkaConsumer
from minio import Minio
from config import config
from pdf2image import convert_from_bytes

from torchvision.utils import make_grid
from torchvision.transforms import transforms
from status_notification_service import notify_status

import os
import logging
import tempfile
import constants
import torchvision.transforms.v2

consumer = KafkaConsumer(
    config["kafka"]["topic"],
    bootstrap_servers=config["kafka"]["brokerURL"],
    group_id=config["kafka"]["groupId"],
    auto_offset_reset="earliest",
)

access_key = os.getenv("MINIO_ACCESS_KEY", "minio")
access_secret = os.getenv("MINIO_SECRET_KEY", "password")
if access_key is None or access_secret is None:
    logging.error("Minio credentials not set")
    exit(1)

minio_client = Minio(
    config["minioEndpoint"],
    access_key=access_key,
    secret_key=access_secret,
    secure=False,
)


def convert(limit_messages=None):
    message_counter = 0
    for message in consumer:
        message_counter += 1
        logging.info("Received message from Kafka: %s", message.value)
        parsed_message = json.loads(message.value.decode("utf-8"))

        checksum = parsed_message["checksum"]
        logging.info("Document checksum: %s", checksum)

        response = None
        try:
            bucket = config["minioBucket"]

            if not minio_client.bucket_exists(bucket):
                logging.error("Bucket %s does not exist", config["minioBucket"])
                return

            logging.info("Retrieving document with checksum %s", checksum)
            response = minio_client.get_object(bucket, f"{checksum}/{constants.default_document_name}")
            document = response.read()

            logging.info("Converting document with checksum %s to images", checksum)
            images = convert_from_bytes(document, 200, fmt="jpeg", output_folder=None)

            logging.info("Uploading images to Minio")
            store_pdf_as_image(bucket, checksum, images)

            logging.info("Successfully uploaded %d images to Minio", len(images))
            logging.info("Document with checksum %s converted successfully", checksum)
            consumer.commit()
            notify_status(checksum, "success")
            if limit_messages is not None and message_counter >= limit_messages:
                logging.info("Limit of %d messages reached", limit_messages)
                consumer.commit()
                consumer.close()
                return
        except Exception as e:
            logging.error("Error during PDF conversion: %s", e)
            notify_status(checksum, "error")
            consumer.commit()
            consumer.close()
            raise e
        finally:
            if response is None:
                return
            else:
                response.close()
                response.release_conn()


def store_pdf_as_image(bucket, checksum, images):
    transform = transforms.Compose([transforms.PILToTensor()])
    tensors = []

    logging.info("Generating image grid")
    for i, image in enumerate(images):
        tensors.append(transform(image))

    with tempfile.TemporaryDirectory() as tmpdirname:
        file_name = f"{tmpdirname}/{constants.image_file_name}"

        logging.info("Saving image grid to %s", file_name)
        grid = make_grid(tensors, nrow=4, padding=5)
        torchvision.transforms.ToPILImage()(grid).save(file_name)

        logging.info("Uploading image grid to Minio")
        minio_client.fput_object(
            bucket,
            f"{checksum}/images/{constants.image_file_name}",
            f"{tmpdirname}/{constants.image_file_name}",
        )
