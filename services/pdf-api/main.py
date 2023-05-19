import logging
from convert_pdf_to_image import convert


def main():
    logging.basicConfig(level=logging.INFO, format="%(asctime)s - %(levelname)s : %(message)s")
    logging.info("Initializing PDF API modules")

    try:
        convert()
    except Exception as e:
        logging.error("Error during PDF conversion: %s", e)
        raise e


main()
