package models

const (
	GameLevelLow      GameLevel = "low"
	GameLevelLowPlus  GameLevel = "low_plus"
	GameLevelMidMinus GameLevel = "mid_minus"
	GameLevelMid      GameLevel = "mid"
	GameLevelMidPlus  GameLevel = "mid_plus"
	GameLevelHigh     GameLevel = "high"
	GameLevelHighPlus GameLevel = "high+"
)

var enToRuGameLevel = map[GameLevel]string{ //nolint:gochecknoglobals
	GameLevelLow:      "начальный",
	GameLevelLowPlus:  "начальный плюс",
	GameLevelMidMinus: "средний минус",
	GameLevelMid:      "средний",
	GameLevelMidPlus:  "средний плюс",
	GameLevelHigh:     "полу-профи",
	GameLevelHighPlus: "профи",
}

func EnToRuGameLevel(level GameLevel) (string, bool) {
	result, ok := enToRuGameLevel[level]
	return result, ok
}

var ruToEnGameLevel = map[string]GameLevel{ //nolint:gochecknoglobals
	"начальный":      GameLevelLow,
	"начальный плюс": GameLevelLowPlus,
	"средний минус":  GameLevelMidMinus,
	"средний":        GameLevelMid,
	"средний плюс":   GameLevelMidPlus,
	"полу-профи":     GameLevelHigh,
	"профи":          GameLevelHighPlus,
}

func RuToEnGameLevel(level string) (GameLevel, bool) {
	result, ok := ruToEnGameLevel[level]
	return result, ok
}

type GameLevel string
