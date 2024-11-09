#!/usr/bin/env python
# pylint: disable=unused-argument
# This program is dedicated to the public domain under the CC0 license.

"""
Basic example for a bot that uses inline keyboards. For an in-depth explanation, check out
 https://github.com/python-telegram-bot/python-telegram-bot/wiki/InlineKeyboard-Example.
"""
import asyncio
import logging
import os
import time

from aiohttp import web
from models import EventCreatedRequest
from telegram import InlineKeyboardButton, InlineKeyboardMarkup, Update, WebAppInfo
from telegram.ext import (
    Application,
    CallbackQueryHandler,
    CommandHandler,
    ContextTypes,
    Updater,
)

# Enable logging
logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s", level=logging.INFO
)
# set higher logging level for httpx to avoid all GET and POST requests being logged
logging.getLogger("httpx").setLevel(logging.WARNING)

LOGGER = logging.getLogger(__name__)

api_app = web.Application()
bot_application = Application.builder().token(os.getenv("BOT_TOKEN")).build()


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Sends a message with three inline buttons attached."""
    chat_id = update.message.chat_id
    # TODO: add web app info
    keyboard = [
        [
            InlineKeyboardButton(
                "Главная",
                url="https://t.me/ond_sportify_test_bot?startapp=main",
                # web_app=WebAppInfo(url="https://91.219.227.107"),
            ),
            InlineKeyboardButton(
                "Создать событие",
                url=f"https://t.me/ond_sportify_test_bot?startapp=create_event__{chat_id}",
            ),
        ],
        [
            InlineKeyboardButton(
                "Карта",
                url="https://t.me/ond_sportify_test_bot?startapp=map",
            ),
        ],
    ]

    reply_markup = InlineKeyboardMarkup(keyboard)

    await update.message.reply_text(
        "Выберете действие",
        reply_markup=reply_markup,
    )


async def help_command(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Displays info on how to use the bot."""
    LOGGER.debug("help_command")
    await update.message.reply_text("Use /start to test this bot.")


async def handle_event_created(request: web.Request) -> web.Response:
    try:
        try:
            data = await request.json()
        except ValueError:
            LOGGER.exception("Error parsing request data")
            return web.json_response(
                {"status": "fail", "reason": "Invalid JSON"},
                status=400,
            )
        try:
            erm = EventCreatedRequest.from_dict(data)
        except TypeError as e:
            LOGGER.exception(f"Error parsing request data {data}")
            return web.json_response({"status": "fail", "reason": str(e)}, status=400)

        message = str(erm.event)

        try:
            await bot_application.bot.send_message(chat_id=erm.tg_chat_id, text=message)
        except Exception as e:
            LOGGER.exception(f"Error sending message to chat {erm.tg_chat_id}")
            return web.json_response({"status": "fail", "reason": str(e)}, status=500)

        return web.json_response({"status": "success"})
    except Exception as e:
        print(f"Error handling event creation: {e}")
        return web.json_response(
            {"status": "fail", "reason": "Internal server error"}, status=500
        )


api_app.router.add_post("/event/created", handle_event_created)


def main() -> None:
    """Run the bot."""
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    bot_application.add_handler(CommandHandler("start", start))
    # application.add_handler(CallbackQueryHandler(button))
    bot_application.add_handler(CommandHandler("help", help_command))

    loop.run_until_complete(bot_application.initialize())
    # await application.post_init()
    loop.run_until_complete(bot_application.updater.start_polling())
    loop.run_until_complete(bot_application.start())

    runner = web.AppRunner(api_app)
    loop.run_until_complete(runner.setup())
    site = web.TCPSite(runner, "0.0.0.0", 8081)
    loop.run_until_complete(site.start())
    print(f"HTTP server started at http://0.0.0.0:8081")

    try:
        loop.run_forever()
    except KeyboardInterrupt:
        LOGGER.info("Received exit signal")
        loop.run_until_complete(bot_application.updater.stop())
        loop.run_until_complete(bot_application.stop())
        loop.run_until_complete(bot_application.shutdown())


if __name__ == "__main__":
    main()
