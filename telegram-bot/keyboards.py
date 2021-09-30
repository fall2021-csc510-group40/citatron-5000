"""
MIT License

Copyright (c) 2021 fall2021-csc510-group40

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
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
