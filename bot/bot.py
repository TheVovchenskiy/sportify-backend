import asyncio
import logging
import os
import time

import httpx
from aiohttp import web
from models.event import (
    EventCreatedRequest,
    EventDeletedRequest,
    EventMessage,
    EventUpdatedRequest,
)
from telegram import InlineKeyboardButton, InlineKeyboardMarkup, Update, WebAppInfo
from telegram.constants import ParseMode
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


event_id_to_message_id: dict[str, EventMessage] = {}


async def subscribe(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Parses the CallbackQuery and updates the message text."""
    query = update.callback_query

    LOGGER.debug(f"Subscribe query: {query}")

    await query.answer()

    target_event_id = None

    for event_id, message in event_id_to_message_id.items():
        if message.message.id == query.message.id:
            target_event_id = event_id
            break

    if target_event_id is None:
        LOGGER.info(f"No event found for message {query.message.id}")
        return

    is_subscribed_response = httpx.get(
        f"http://0.0.0.0:8090/api/v1/events/{target_event_id}/subscribers?tg_id={query.from_user.id}",
    )

    LOGGER.info(f"Is subscribed response: {is_subscribed_response.json()}")

    if "is_subscribed" not in is_subscribed_response.json():
        LOGGER.error(
            f"Failed to check if user {query.from_user.id} is subscribed to event {target_event_id}, error: {is_subscribed_response.text}"
        )
        return

    resp = httpx.put(
        f"http://0.0.0.0:8090/api/v1/events/{target_event_id}/subscribers",
        json={
            "sub": not is_subscribed_response.json()["is_subscribed"],
            "tg_id": query.from_user.id,
        },
    )

    LOGGER.info(f"Subscribe response: {resp.json()}")

    if 200 <= resp.status_code < 300:
        LOGGER.info(f"Successfully subscribed/unsubscribed to event {target_event_id}")
    else:
        LOGGER.error(
            f"Failed to subscribe/unsubscribe to event {target_event_id}, error: {resp.text}"
        )


async def start(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Sends a message with three inline buttons attached."""
    command: str = update.message.text

    # tokens as tokens of the command
    tokens = command.split()

    if len(tokens) > 1:
        # token as 2nd argument of the command, should be token passed from
        # t.me/<bot_name>?start=<token>
        token = tokens[1]
        user_id = update.message.from_user.id
        chat_id = update.message.chat_id
        username = update.message.from_user.username
        first_name = update.message.from_user.first_name
        chat_type = update.message.chat.type
        # TODO: handle invalid tokens !!!
        LOGGER.info(
            "Handling start command with token, ("
            f"user_id = {user_id}"
            f"username = {username}"
        )

        resp = httpx.post(
            f"http://0.0.0.0:8090/api/v1/users",
            json={
                "tg_user_id": user_id,
                "tg_username": username,
                "token": token,
                "tg_update": {
                    "update_id": update.update_id,
                    "message": {
                        "chat": {
                            "id": chat_id,
                            "first_name": username,
                            "type": chat_type,
                        },
                        "text": command,
                    },
                },
            },
        )
        LOGGER.info(f"User creation response status code: {resp.status_code}")
        if 200 <= resp.status_code < 300:
            LOGGER.info(f"Successfully authenticated user {user_id}")
            await update.message.reply_text(
                f"✅ Вы успешно вошли, вернитесь пожалуйста обратно на сайт"
            )
        else:
            LOGGER.error(f"Failed to authenticate user {user_id}, error: {resp.text}")
            await update.message.reply_text(
                "❌ Произошла ошибка при авторизации, попробуйте еще раз"
            )
        return

    chat_id = update.message.chat_id
    # TODO: add web app info
    keyboard = [
        [
            InlineKeyboardButton(
                "Главная",
                url=f"https://t.me/ond_sportify_bot?startapp=events__{chat_id}",
            ),
            InlineKeyboardButton(
                "Создать событие",
                url=f"https://t.me/ond_sportify_bot?startapp=create_event__{chat_id}",
            ),
        ],
        [
            InlineKeyboardButton(
                "Карта",
                url=f"https://t.me/ond_sportify_bot?startapp=map__{chat_id}",
            ),
        ],
    ]

    reply_markup = InlineKeyboardMarkup(keyboard)

    await update.message.reply_text(
        "Выберете действие",
        reply_markup=reply_markup,
    )


async def test(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Sends a message with three inline buttons attached."""
    keyboard = [
        [
            InlineKeyboardButton(
                "Записаться/Отписаться",
                callback_data="1",
            ),
        ],
        [],
    ]

    reply_markup = InlineKeyboardMarkup(keyboard)

    await update.message.reply_text(
        "Выберете действие",
        reply_markup=reply_markup,
    )


async def help_command(update: Update, context: ContextTypes.DEFAULT_TYPE) -> None:
    """Displays info on how to use the bot."""
    LOGGER.debug(f"Received /help command (update={update})")
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

        try:
            message = str(erm.event)
        except Exception as e:
            LOGGER.exception(f"Error creating message for event {repr(erm)}")
            return web.json_response({"status": "fail", "reason": str(e)}, status=500)

        try:
            keyboard = [
                [
                    InlineKeyboardButton(
                        "Записаться/Отписаться",
                        callback_data="1",
                    ),
                ],
            ]

            LOGGER.info(f"Sending message to chat {erm.tg_chat_id}")
            message = await bot_application.bot.send_photo(
                chat_id=erm.tg_chat_id,
                photo=erm.event.url_preview,
                caption=message,
                parse_mode=ParseMode.MARKDOWN_V2,
                reply_markup=InlineKeyboardMarkup(keyboard),
            )
            event_id_to_message_id[erm.event.id] = EventMessage(
                event=erm.event,
                message=message,
            )
        except Exception as e:
            LOGGER.exception(f"Error sending message to chat {erm.tg_chat_id}")
            return web.json_response({"status": "fail", "reason": str(e)}, status=500)

        return web.json_response({"status": "success"})
    except Exception as e:
        LOGGER.exception(f"Error handling event creation")
        return web.json_response(
            {"status": "fail", "reason": "Internal server error"}, status=500
        )


async def handle_event_updated(request: web.Request) -> web.Response:
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
            erm = EventUpdatedRequest.from_dict(data=data)
        except TypeError as e:
            LOGGER.exception(f"Error parsing request data {data}")
            return web.json_response({"status": "fail", "reason": str(e)}, status=400)

        if erm.event.id not in event_id_to_message_id:
            LOGGER.info(f"Event {erm.event.id} not found")
            return web.json_response(
                {"status": "fail", "reason": "Event not found"}, status=404
            )

        erm_old = event_id_to_message_id[erm.event.id]

        try:
            new_message = str(erm.event)
        except Exception as e:
            LOGGER.exception(f"Error creating message for event {repr(erm)}")
            return web.json_response({"status": "fail", "reason": str(e)}, status=500)

        try:
            keyboard = [
                [
                    InlineKeyboardButton(
                        "Записаться/Отписаться",
                        callback_data="1",
                    ),
                ],
            ]

            await bot_application.bot.edit_message_caption(
                chat_id=erm_old.message.chat.id,
                caption=new_message,
                message_id=erm_old.message.id,
                parse_mode=ParseMode.MARKDOWN_V2,
                reply_markup=InlineKeyboardMarkup(keyboard),
            )
        except Exception as e:
            LOGGER.exception(f"Error editing event {erm.event.id!r}")
            return web.json_response({"status": "fail", "reason": str(e)}, status=500)

        return web.json_response({"status": "success"})
    except Exception as e:
        LOGGER.exception(f"Error handling event creation")
        return web.json_response(
            {"status": "fail", "reason": "Internal server error"}, status=500
        )


async def handle_event_deleted(request: web.Request) -> web.Response:
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
            erm = EventDeletedRequest.from_dict(data=data)
        except TypeError as e:
            LOGGER.exception(f"Error parsing request data {data}")
            return web.json_response({"status": "fail", "reason": str(e)}, status=400)

        if erm.event_id not in event_id_to_message_id:
            LOGGER.info(f"Event {erm.event_id} not found")
            return web.json_response(
                {"status": "fail", "reason": "Event not found"}, status=404
            )

        erm_old = event_id_to_message_id[erm.event_id]

        try:
            await bot_application.bot.delete_message(
                chat_id=erm_old.message.chat.id,
                message_id=erm_old.message.id,
            )
        except Exception as e:
            LOGGER.exception(f"Error deleting event {erm.event_id!r}")
            return web.json_response({"status": "fail", "reason": str(e)}, status=500)

        return web.json_response({"status": "success"})
    except Exception as e:
        LOGGER.exception(f"Error handling event creation")
        return web.json_response(
            {"status": "fail", "reason": "Internal server error"}, status=500
        )


api_app.router.add_post("/event/created", handle_event_created)
api_app.router.add_put("/event/updated", handle_event_updated)
api_app.router.add_delete("/event/deleted", handle_event_deleted)


def main() -> None:
    """Run the bot."""
    loop = asyncio.new_event_loop()
    asyncio.set_event_loop(loop)

    bot_application.add_handler(CommandHandler("start", start))
    bot_application.add_handler(CommandHandler("test", test))
    bot_application.add_handler(CallbackQueryHandler(subscribe))
    bot_application.add_handler(CommandHandler("help", help_command))

    loop.run_until_complete(bot_application.initialize())
    # await application.post_init()
    loop.run_until_complete(
        bot_application.updater.start_polling(allowed_updates=Update.ALL_TYPES)
    )
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
