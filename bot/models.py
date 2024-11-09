from dataclasses import dataclass


@dataclass
class Event:
    id: str
    creator_id: str
    sport_type: str
    address: str
    date: str
    start_time: str
    end_time: str
    price: int
    is_free: bool
    game_level: list[str]
    capacity: int
    busy: int
    subscribers_id: list[str]
    preview: str
    photos: list[str]
    latitude: str
    longitude: str

    def __str__(self) -> str:
        lines = [
            "Событие:",
            f"ID: {self.id}",
            f"Создатель: {self.creator_id}",
            f"Вид спорта: {self.sport_type}",
            f"Адрес: {self.address}",
            f"Дата: {self.date}",
            f"Время начала: {self.start_time}",
            f"Время окончания: {self.end_time}",
            f"Цена: {self.price}",
            f"Бесплатно: {'Да' if self.is_free else 'Нет'}",
            f"Уровень игры: {', '.join(self.game_level)}",
            f"Вместимость: {self.capacity}",
            f"Занято мест: {self.busy}",
            f"Подписчики: {', '.join(self.subscribers_id)}",
            f"Превью: {self.preview}",
        ]

        return "\n".join(lines)


@dataclass
class EventCreatedRequest:
    tg_chat_id: str
    tg_user_id: str
    event: Event

    @classmethod
    def from_dict(cls, data: dict):
        event_data = data.pop("event")
        event = Event(**event_data)
        return cls(event=event, **data)


if __name__ == "__main__":
    data = {
        "tg_chat_id": "123",
        "tg_user_id": "456",
        "invalid": 456,
    }

    erm = EventCreatedRequest(**data)
