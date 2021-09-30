"""
Copyright (c) 2021 contributors of the Citatron-5000 Project. All Rights Reserved
"""
import json

from telegram.ext import Updater, MessageHandler, CallbackQueryHandler, CallbackContext, Filters, ConversationHandler
from telegram import Update
import argparse
import logging

from keyboards import get_choice_keyboard
from api import APIException, APIAdaptor

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s',
                    level=logging.INFO)


def handle_errors(handler):
    def wrapped(self, update, context):
        try:
            return handler(self, update, context)
        except APIException as e:
            error_string = f'Call to handle {e.handle} resulted in error, status_code: {e.status_code}'
            if e.error:
                error_string = error_string + f", error_text: {e.error}"
            context.bot.send_message(chat_id=update.effective_chat.id, text=error_string)
            return ConversationHandler.END

    return wrapped


class Bot:
    CHOOSING, FORMATTING = range(2)
    SEARCH_PAGE_SIZE = 5

    def __init__(self, config_path):
        self.config = self._load_config(config_path)
        self.logger = logging.getLogger("citebot")
        self.logger.setLevel(self.config["log_level"])
        self.api = APIAdaptor(self.config["api_url"], self.logger)

    @staticmethod
    def _load_config(config_path):
        with open(config_path, "r") as f:
            return json.load(f)

    def _request_search_choice(self, update, context):
        results = context.chat_data["results"]
        offset = context.chat_data["results_offset"]
        assert offset < len(results)

        has_next_page = (len(results) - offset > Bot.SEARCH_PAGE_SIZE)
        has_prev_page = (offset > 0)
        self.logger.debug(f"creating choice keyboard, total: {len(results)}, offset: {offset}, next: {has_next_page}, prev: {has_prev_page}")
        keyboard = get_choice_keyboard([
            {
                "text": f'{work["title"]} by {", ".join(work["authors"])}',
                "data": work["id"]
            } for work in results[offset:offset + Bot.SEARCH_PAGE_SIZE]
        ], add_next_page=has_next_page, add_prev_page=has_prev_page)
        self.logger.debug(f"Sending keyboard: {keyboard.inline_keyboard}")

        context.bot.send_message(text="Please choose the desired work:",
                                 chat_id=update.effective_chat.id,
                                 reply_markup=keyboard)

    @handle_errors
    def search(self, update: Update, context: CallbackContext):
        title = update.message.text
        results = self.api.get_search_results(title)

        self.logger.debug(f"Got search results: {results}")

        if not results:
            self.logger.debug("No results, stopping conversation")
            context.bot.send_message(text="No entries were found =(",
                                     chat_id=update.effective_chat.id)
            return ConversationHandler.END

        if len(results) == 1:
            self.logger.info(f"Single results found, using id {results[0]['id']}")
            context.chat_data["id_to_format"] = results[0]["id"]
            return self.choose_format(update, context)

        context.chat_data["results"] = results
        context.chat_data["results_offset"] = 0
        self._request_search_choice(update, context)
        return Bot.CHOOSING

    def choose_format(self, update: Update, context: CallbackContext):
        self.logger.info("Choosing format")
        keyboard = get_choice_keyboard([
            {
                "text": description,
                "data": format_name
            } for format_name, description in self.config["formats"].items()
        ])

        context.bot.send_message(text="Please choose format:",
                                 chat_id=update.effective_chat.id,
                                 reply_markup=keyboard)

        return Bot.FORMATTING

    def handle_search_choice(self, update: Update, context: CallbackContext):
        self.logger.info("Handling search choice")
        query = update.callback_query

        context.bot.delete_message(chat_id=update.effective_chat.id,
                                   message_id=update.effective_message.message_id)

        if query.data == "cancel":
            context.bot.send_message(text="Search canceled",
                                     chat_id=update.effective_chat.id)
            return ConversationHandler.END
        elif query.data == "next_page":
            context.chat_data["results_offset"] += Bot.SEARCH_PAGE_SIZE
            self._request_search_choice(update, context)
            return Bot.CHOOSING
        elif query.data == "prev_page":
            context.chat_data["results_offset"] -= Bot.SEARCH_PAGE_SIZE
            self._request_search_choice(update, context)
            return Bot.CHOOSING
        else:
            context.chat_data["id_to_format"] = query.data
            return self.choose_format(update, context)

    @handle_errors
    def handle_format_choice(self, update: Update, context: CallbackContext):
        query = update.callback_query
        context.chat_data["format_name"] = query.data

        context.bot.delete_message(chat_id=update.effective_chat.id,
                                   message_id=update.effective_message.message_id)

        resulting_text = self.api.get_format_results(context.chat_data["id_to_format"],
                                                     context.chat_data["format_name"])
        context.bot.send_message(text=f"```\n{resulting_text}\n```",
                                 parse_mode="MarkdownV2",
                                 chat_id=update.effective_chat.id)

        return ConversationHandler.END

    def run(self):
        updater = Updater(token=self.config["token"])
        dispatcher = updater.dispatcher

        dispatcher.add_handler(ConversationHandler(
            entry_points=[MessageHandler(callback=self.search, filters=Filters.text & ~Filters.command)],
            per_message=False,
            states={
                Bot.CHOOSING: [
                    CallbackQueryHandler(self.handle_search_choice)
                ],
                Bot.FORMATTING: [
                    CallbackQueryHandler(self.handle_format_choice)
                ]
            },
            fallbacks=[],
        ))

        updater.start_polling()
        updater.idle()


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("config", metavar="CONFIG_PATH", type=str, nargs=1, help="Config path")

    args = parser.parse_args()

    bot = Bot(args.config[0])
    bot.run()


if __name__ == "__main__":
    main()
