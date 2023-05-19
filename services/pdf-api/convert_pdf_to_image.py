# create a kafka consumer
import json

from kafka import KafkaConsumer
from minio import Minio
from config import config
from pdf2image import convert_from_bytes

import os
import logging
import tempfile

consumer = KafkaConsumer(
    config["kafka"]["topic"],
    bootstrap_servers=config["kafka"]["brokerURL"],
    group_id=config["kafka"]["groupId"],
    auto_offset_reset="earliest",
)

access_key = os.getenv("MINIO_ACCESS_KEY")
access_secret = os.getenv("MINIO_SECRET_KEY")
if access_key is None or access_secret is None:
    logging.error("Minio credentials not set")
    exit(1)

minio_client = Minio(
    config["minioEndpoint"],
    access_key=access_key,
    secret_key=access_secret,
    secure=False,
)


def convert():
    for message in consumer:
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
            response = minio_client.get_object(bucket, "{}/document.pdf".format(checksum))
            document = response.read()

            logging.info("Converting document with checksum %s to images", checksum)
            images = convert_from_bytes(document, 200, fmt="jpeg", output_folder=None)

            logging.info("Uploading images to Minio")
            for i, image in enumerate(images):
                with tempfile.TemporaryDirectory() as tmpdirname:
                    image.save("{}/page-{}.jpg".format(tmpdirname, i), "JPEG")
                    minio_client.fput_object(
                        bucket,
                        "{}/images/page-{}.jpg".format(checksum, i),
                        "{}/page-{}.jpg".format(tmpdirname, i),
                    )

            logging.info("Successfully uploaded %d images to Minio", len(images))
            logging.info("Document with checksum %s converted successfully", checksum)
            consumer.commit()
        except Exception as e:
            logging.error("Error during PDF conversion: %s", e)
            consumer.commit()
            consumer.close()
            raise e
        finally:
            if response is None:
                return
            else:
                response.close()
                response.release_conn()
