from dataclasses import dataclass

from telegram.helpers import escape_markdown


@dataclass
class User:
    id: str
    username: str
    tg_id: int | None = None

    @classmethod
    def from_dict(cls, data: dict) -> "User":
        return cls(**data)

    def __str__(self) -> str:
        return (
            f"[{escape_markdown(self.username, 2)}](tg://user?id={self.tg_id})"
            if self.tg_id
            else escape_markdown(self.username, 2)
        )
