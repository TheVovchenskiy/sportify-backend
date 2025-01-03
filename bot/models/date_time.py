import datetime
from dataclasses import dataclass

from telegram.helpers import escape_markdown

DATETIME_FORMAT = "%Y-%m-%dT%H:%M:%SZ"
DATE_FORMAT = "%d.%m.%Y"
TIME_FORMAT = "%H:%M"


@dataclass
class DateTime:
    date: datetime.datetime
    start_time: datetime.datetime
    end_time: datetime.datetime | None = None

    @classmethod
    def from_dict(cls, data: dict) -> "DateTime":
        return cls(
            date=datetime.datetime.strptime(data["date"], DATETIME_FORMAT),
            start_time=datetime.datetime.strptime(data["start_time"], DATETIME_FORMAT),
            end_time=(
                datetime.datetime.strptime(data["end_time"], DATETIME_FORMAT)
                if data.get("end_time")
                else None
            ),
        )

    def __str__(self) -> str:
        parts = [
            f"📅 *Дата*: {escape_markdown(self.date.strftime(DATE_FORMAT), 2)}",
            f"🕘 *Начало*: {escape_markdown(self.start_time.strftime(TIME_FORMAT),2)}",
        ]
        if self.end_time:
            parts.append(
                f"🕥 *Конец*: {escape_markdown(self.end_time.strftime(TIME_FORMAT), 2)}",
            )

        return "\n".join(parts)


if __name__ == "__main__":
    data = {
        "date": "2024-11-24T00:00:00Z",
        "start_time": "2024-01-01T00:04:00Z",
        "end_time": "2024-01-01T00:04:05Z",
    }

    xt = DateTime.from_dict(data)
    print(xt)
