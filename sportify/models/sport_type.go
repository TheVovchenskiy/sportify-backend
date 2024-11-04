package models

const (
	SportTypeVolleyball SportType = "volleyball"
	SportTypeBasketball SportType = "basketball"
	SportTypeFootball   SportType = "football"
)

type SportType string

var enToRuSportType = map[SportType]string{ //nolint:gochecknoglobals
	SportTypeVolleyball: "волейбол",
	SportTypeBasketball: "баскетбол",
	SportTypeFootball:   "футбол",
}

func EnToRuSportType(sportType SportType) (string, bool) {
	result, ok := enToRuSportType[sportType]
	return result, ok
}

var ruToEnSportType = map[string]SportType{ //nolint:gochecknoglobals
	"волейбол":  SportTypeVolleyball,
	"баскетбол": SportTypeBasketball,
	"футбол":    SportTypeFootball,
}

func RuToEnSportType(sportType string) (SportType, bool) {
	result, ok := ruToEnSportType[sportType]
	return result, ok
}
