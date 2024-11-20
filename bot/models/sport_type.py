from enum import Enum


class SportType(Enum):
    FOOTBALL = "football"
    BASKETBALL = "basketball"
    VOLLEYBALL = "volleyball"


en_to_ru_sport_type = {
    SportType.FOOTBALL: "футбол",
    SportType.BASKETBALL: "баскетбол",
    SportType.VOLLEYBALL: "волейбол",
}


def get_sport_type_ru(sport_type: str) -> str:
    # TODO: check if present in map
    return en_to_ru_sport_type[SportType(sport_type)]
