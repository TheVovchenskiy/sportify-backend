from enum import Enum


class GameLevel(Enum):
    LOW = "low"
    LOW_PLUS = "low_plus"
    MID_MINUS = "mid_minus"
    MID = "mid"
    MID_PLUS = "mid_plus"
    HIGH = "high"
    HIGH_PLUS = "high_plus"


en_to_ru_game_level = {
    GameLevel.LOW: "начальный",
    GameLevel.LOW_PLUS: "начальный+",
    GameLevel.MID_MINUS: "средний-",
    GameLevel.MID: "средний",
    GameLevel.MID_PLUS: "средний+",
    GameLevel.HIGH: "полу-профи",
    GameLevel.HIGH_PLUS: "профи",
}

def get_game_level_ru(sport_type: str) -> str:
    # TODO: check if present in map
    return en_to_ru_game_level[GameLevel(sport_type)]
