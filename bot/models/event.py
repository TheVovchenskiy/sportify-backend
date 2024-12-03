from dataclasses import dataclass

from models.date_time import DateTime
from models.game_level import GameLevel, en_to_ru_game_level
from models.sport_type import get_sport_type_ru
from models.user import User
from telegram import Message
from telegram.helpers import escape_markdown


@dataclass
class Event:
    id: str
    creator: User
    sport_type: str
    address: str
    date_and_time: DateTime
    is_free: bool
    game_levels: list[str]
    busy: int
    subscribers: list[User]
    url_preview: str | None = None
    capacity: int | None = None
    price: int | None = None
    latitude: str | None = None
    longitude: str | None = None
    hashtags: list[str] | None = None

    def __str__(self) -> str:
        lines = [
            "*Событие*",
            f"Автор: {self.creator}",
            f"Вид спорта: {get_sport_type_ru(self.sport_type)}",
            escape_markdown(f"Адрес: {self.address}", 2),
            escape_markdown(str(self.date_and_time), 2),
            f"Цена: {self.price if not self.is_free else "БЕСПЛАТНО"}",
            (
                f"Уровень игры: [{', '.join(f"`{escape_markdown(en_to_ru_game_level[GameLevel(game_level)], 2) }`" for game_level in self.game_levels)}]"
                if self.game_levels
                else ""
            ),
            f"Вместимость: {self.capacity}",
            f"Занято мест: {self.busy}",
            (f"Участники:\n{self.__str_subscribers()}" if self.subscribers else ""),
        ]

        if self.hashtags:
            lines.append("\n")
            lines.append(" ".join(self.hashtags))

        return "\n".join(lines)

    def __str_subscribers(self):
        if self.subscribers:
            return "\n".join(
                escape_markdown("- ", 2) + str(subscriber)
                for subscriber in self.subscribers
            )

        return ""

    @classmethod
    def from_dict(cls, data: dict) -> "Event":
        creator_data = data.pop("creator")
        creator = User.from_dict(creator_data)

        subscribers_data = data.pop("subscribers")
        subscribers = [
            User.from_dict(subscriber_data) for subscriber_data in subscribers_data
        ]

        date_and_time_data = data.pop("date_and_time")
        date_and_time = DateTime.from_dict(**date_and_time_data)

        return cls(
            creator=creator,
            subscribers=subscribers,
            date_and_time=date_and_time,
            **data,
        )


@dataclass
class EventCreatedRequest:
    tg_chat_id: str
    # tg_user_id: str
    event: Event

    @classmethod
    def from_dict(cls, data: dict):
        event_data = data.pop("event")
        event = Event.from_dict(event_data)
        return cls(event=event, **data)


@dataclass
class EventUpdatedRequest:
    event: Event

    @classmethod
    def from_dict(cls, data: dict):
        event_data = data.pop("event")
        event = Event.from_dict(event_data)
        return cls(event=event, **data)


@dataclass
class EventDeletedRequest:
    event_id: str

    @classmethod
    def from_dict(cls, data: dict):
        return cls(**data)


@dataclass
class EventMessage:
    event: Event
    message: Message
    # chat_id: str
    # message_id: int
