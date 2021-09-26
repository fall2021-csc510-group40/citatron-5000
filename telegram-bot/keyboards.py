from telegram import InlineKeyboardButton, InlineKeyboardMarkup
import emoji


def _get_digit_emoji(number):
    return f"{number}\U0000FE0F\U000020E3"


def get_choice_keyboard(choices, add_next_page=False):
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

    if add_next_page:
        keyboard_layout += [[
            InlineKeyboardButton(
                text=emoji.emojize(":magnifying_glass_tilted_right:", use_aliases=True) + " - " + "Next page...",
                callback_data="next_page"
            )
        ]]

    keyboard_layout += [[
        InlineKeyboardButton(
            text=emoji.emojize(":cross_mark:", use_aliases=True) + " - " + "Cancel",
            callback_data="cancel"
        )
    ]]

    keyboard = InlineKeyboardMarkup(keyboard_layout)

    return keyboard
