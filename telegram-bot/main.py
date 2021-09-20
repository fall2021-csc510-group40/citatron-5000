import json

import requests
from telegram.ext import Updater, CommandHandler, MessageHandler, CallbackQueryHandler, CallbackContext, Filters
from telegram import Update
import argparse
import logging


logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                     level=logging.INFO)


API_URL = "127.0.0.1:8899"


class APIException(Exception):
    def __init__(self, handle, status_code, error):
        self.handle = handle
        self.status_code = status_code
        self.error = error


def handle_errors(handler):
    def wrapped(update, context):
        try:
            handler(update, context)
        except APIException as e:
            error_string = f'Call to handle {e.handle} resulted in error, status_code: {e.status_code}'
            if e.error:
                error_string = error_string + f", error_text: {e.error}"
            context.bot.send_message(chat_id=update.effective_chat.id, text=error_string)

    return wrapped


def get_search_results(title):
    rsp = requests.get(API_URL + "/search", data={
        "query": {
            "title": title,
        },
    })

    if not rsp.ok:
        raise APIException("search", rsp.status_code, rsp.json().get("error", None))

    return rsp.json()["results"]


def get_format_results(item_id, format):
    rsp = requests.get(API_URL + "/format", data={
        "id": item_id,
        "format": format,
    })

    if not rsp.ok:
        raise APIException("format", rsp.status_code, rsp.json().get("error", None))

    return rsp.json()["text"]


def button(update: Update, context: CallbackContext):
    query = update.callback_query
    query.answer()
    context.bot.send_message(text=str(query.data))


def start(update, context):

    logging.info(f"Request received, chat_id: {update.effective_chat.id}, username: {update.effective_chat.username}")


@handle_errors
def search(update: Update, context: CallbackContext):
    title = update.message.text
    results = get_search_results(title)
    context.chat_data["state"] = "choose"
    context.chat_data["results"] = results




def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("config", metavar="CONFIG_PATH", type=str, nargs=1, help="Config path")

    args = parser.parse_args()

    with open(args.config[0], "r") as f:
        config = json.load(f)

    updater = Updater(token=config["token"])
    dispatcher = updater.dispatcher
    start_handler = CommandHandler("start", start)
    dispatcher.add_handler(start_handler)
    dispatcher.add_handler(CallbackQueryHandler(button))

    updater.start_polling()


if __name__ == "__main__":
    main()
