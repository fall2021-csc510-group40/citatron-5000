from telegram import InlineKeyboardButton, InlineKeyboardMarkup


def _get_digit_emoji(number):
    return f"{number}\U0000FE0F\U000020E3"


def get_choice_keyboard(choices):
    if len(choices) > 9:
        raise ValueError()

    keyboard = InlineKeyboardMarkup([
        [InlineKeyboardButton(
            text=_get_digit_emoji(i + 1) + " - " + option["text"],
            callback_data=option["data"])
        ] for i, option in enumerate(choices)
    ])

    return keyboard
