from dataclasses import dataclass


@dataclass
class DateTime:
    date: str
    start_time: str
    end_time: str | None = None
