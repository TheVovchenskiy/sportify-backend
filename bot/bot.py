#!/usr/bin/env python
import logging

import httpx
from telegram import Message, Update
from telegram.ext import Application, ContextTypes, MessageHandler, filters

from config import BotConfig

logging.basicConfig(
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s",
    level=logging.INFO,
)

# set higher logging level for httpx to avoid all GET and POST requests being logged
logging.getLogger("httpx").setLevel(logging.WARNING)

logger = logging.getLogger(__name__)


class Bot:
    """
    Class representing a Telegram bot for parsing chats and messages.
    """

    config: BotConfig

    def __init__(self, config: BotConfig):
        self.config = config

        self._application = Application.builder().token(config.token).build()
        self._application.add_handler(MessageHandler(filters.ALL, self._parse_messages))

    async def _parse_messages(self, update: Update, context: ContextTypes.DEFAULT_TYPE):
        message: Message | None = update.message

        if (
            message is None
            or message.text is None  # ignore empty messages
            or message.text.startswith("/")  # ignore commands
            or message.from_user is None  # ignore messages from bots
            or message.chat.title is None  # ignore messages from private chats
        ):
            logger.info(f"Ignoring message: {message}")
            return

        logger.info(
            f"Received message: "
            f"(text = {message.text!r}, "
            f"author = {message.from_user.name!r}, "
            f"chat = {message.chat.title!r})"
        )

        await self._handle_message(message)

    async def _handle_message(self, message: Message):
        """
        Sends a message to the api url.
        """
        logger.info(f"Sending message to {self.config.api_url!r}")
        message_dict = message.to_dict()

        async with httpx.AsyncClient() as client:
            response = await client.post(self.config.api_url, json=message_dict)

        logger.info(f"Sent message to API. Response status: {response.status_code}")

    def run(self):
        """
        Run the bot.
        """
        self._application.run_polling(allowed_updates=[Update.MESSAGE])
