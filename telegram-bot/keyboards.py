"""
Copyright (c) 2021 contributors of the Citatron-5000 Project. All Rights Reserved
"""
from telegram import InlineKeyboardButton, InlineKeyboardMarkup
import emoji


def _get_digit_emoji(number):
    return f"{number}\U0000FE0F\U000020E3"


def get_choice_keyboard(choices, add_next_page=False, add_prev_page=False):
    if len(choices) > 9:
        raise ValueError()

    keyboard_layout = [
        [
            InlineKeyboardButton(
                text=_get_digit_emoji(i + 1) + " - " + option["text"],
                callback_data=option["data"]
            )
        ] for i, option in enumerate(choices)
    ]

    pages = []
    if add_prev_page:
        pages.append(InlineKeyboardButton(
                        text=emoji.emojize(":left_arrow:", use_aliases=True) + " - " + "Previous page...",
                        callback_data="prev_page"
                    ))

    if add_next_page:
        pages.append(InlineKeyboardButton(
                        text=emoji.emojize(":right_arrow:", use_aliases=True) + " - " + "Next page...",
                        callback_data="next_page"
                    ))

    if len(pages) > 0:
        keyboard_layout.append(pages)

    keyboard_layout.append([
        InlineKeyboardButton(
            text=emoji.emojize(":cross_mark:", use_aliases=True) + " - " + "Cancel",
            callback_data="cancel"
        )
    ])

    keyboard = InlineKeyboardMarkup(keyboard_layout)

    return keyboard
