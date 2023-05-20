import unittest

from convert_pdf_to_image import convert
from minio import Minio
from config import config
from kafka import KafkaConsumer, KafkaProducer
from kafka.admin import KafkaAdminClient

import os
import json
import logging


class ConvertTestCase(unittest.TestCase):
    minio_client = None
    kafka_consumer = None
    kafka_producer = None
    checksum = None

    logging.basicConfig(level=logging.INFO, format="TEST %(asctime)s - %(levelname)s : %(message)s")

    def delete_existing_objects(self):
        bucket_objects = self.minio_client.list_objects(config["minioBucket"], recursive=True)
        if len(list(bucket_objects)) > 0:
            objects_to_delete = []
            for bucket_object in bucket_objects:
                objects_to_delete.append(bucket_object.object_name)
            self.minio_client.remove_objects(config["minioBucket"], objects_to_delete)

    # Set up the test environment
    def setUp(self):
        access_key = "minio"
        os.environ["MINIO_ACCESS_KEY"] = access_key

        secret_key = "password"
        os.environ["MINIO_SECRET_KEY"] = secret_key

        self.minio_client = Minio(
            config["minioEndpoint"],
            access_key=access_key,
            secret_key=secret_key,
            secure=False,
        )

        self.kafka_consumer = KafkaConsumer(
            "process-document-status",
            bootstrap_servers="localhost:9092",
            group_id="processed-status-group",
            auto_offset_reset="earliest",
        )

        self.kafka_producer = KafkaProducer(
            bootstrap_servers="localhost:9092",
        )

        # Checksum of the test pdf document
        self.checksum = "bb07b36f090a8721906bed6454f437ceb260a23e36784871ff8bac96979b2c5c"

        # Send event to Kafka
        self.kafka_producer.send("process-document", json.dumps({"checksum": self.checksum}).encode("utf-8"))
        self.kafka_producer.flush()
        self.kafka_producer.close()

        # Remove all objects from the Minio bucket
        if not self.minio_client.bucket_exists(config["minioBucket"]):
            self.minio_client.make_bucket(config["minioBucket"])
        else:
            self.delete_existing_objects()

        self.minio_client.fput_object(config["minioBucket"], "{}/document.pdf".format(self.checksum),
                                      "testdata/document.pdf")

    # Tear down the test environment
    def tearDown(self):
        KafkaAdminClient(bootstrap_servers="localhost:9092").delete_topics(
            ["process-document"])
        self.kafka_consumer.close()
        self.kafka_producer.close()
        self.delete_existing_objects()

    # Test the convert function
    def test_should_convert_pdf_to_images(self):
        minio_client = self.minio_client
        kafka_consumer = self.kafka_consumer

        # Run convert
        convert(limit_messages=1)

        self.assertTrue(minio_client.bucket_exists(config["minioBucket"]))
        objects = minio_client.list_objects(config["minioBucket"], recursive=True)

        # Check that the images were uploaded to Minio
        self.assertEqual(len(list(objects)), 4)

        # Check if status was sent to Kafka
        message = next(kafka_consumer)
        self.assertEqual(json.loads(message.value.decode("utf-8")), {"checksum": self.checksum, "status": "success"})


if __name__ == '__main__':
    unittest.main()
