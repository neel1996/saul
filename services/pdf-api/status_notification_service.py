from kafka import KafkaProducer

from config import config

import json
import logging

producer = KafkaProducer(
    bootstrap_servers=config["kafka"]["brokerURL"],
)


def notify_status(checksum, status):
    logging.info("Notifying status '%s' for document with checksum %s", status, checksum)
    message = {
        "checksum": checksum,
        "status": status
    }

    topic = config["kafka"]["statusTopic"]
    producer.send(topic, json.dumps(message).encode("utf-8"))
    producer.flush()
