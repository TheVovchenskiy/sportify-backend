package models

const (
	SportTypeVolleyball  SportType = "volleyball"
	SportTypeBasketball  SportType = "basketball"
	SportTypeFootball    SportType = "football"
	SportTypeTennis      SportType = "tennis"
	SportTypeTableTennis SportType = "table_tennis"
	SportTypeRunning     SportType = "running"
	SportTypeHockey      SportType = "hockey"
	SportTypeSkating     SportType = "skating"
	SportTypeSkiing      SportType = "skiing"
)

type SportType string

var enToRuSportType = map[SportType]string{ //nolint:gochecknoglobals
	SportTypeVolleyball:  "волейбол",
	SportTypeBasketball:  "баскетбол",
	SportTypeFootball:    "футбол",
	SportTypeTennis:      "теннис",
	SportTypeTableTennis: "настольный теннис",
	SportTypeRunning:     "бег",
	SportTypeHockey:      "хоккей",
	SportTypeSkating:     "катание на коньках",
	SportTypeSkiing:      "катание на лыжах",
}

func EnToRuSportType(sportType SportType) (string, bool) {
	result, ok := enToRuSportType[sportType]
	return result, ok
}

var ruToEnSportType = map[string]SportType{ //nolint:gochecknoglobals
	"волейбол":           SportTypeVolleyball,
	"баскетбол":          SportTypeBasketball,
	"футбол":             SportTypeFootball,
	"теннис":             SportTypeTennis,
	"настольный теннис":  SportTypeTableTennis,
	"бег":                SportTypeRunning,
	"хоккей":             SportTypeHockey,
	"катание на коньках": SportTypeSkating,
	"катание на лыжах":   SportTypeSkiing,
}

func RuToEnSportType(sportType string) (SportType, bool) {
	result, ok := ruToEnSportType[sportType]
	return result, ok
}
