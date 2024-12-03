package app

import (
	"fmt"
	"strings"

	"github.com/TheVovchenskiy/sportify-backend/models"
)

func replaceDashAndSpaceToUnderscore(s string) string {
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, "-", "_")
	return s
}

const formatHashtag = "#%s"

func GenerateHashtags(event *models.ShortEvent) (hashtags []string) {
	hashtags = append(hashtags, event.DateAndTime.Date.Format("#дата_02_01_2006"))

	if sportType, ok := models.EnToRuSportType(event.SportType); ok {
		hashtags = append(hashtags, fmt.Sprintf(formatHashtag, replaceDashAndSpaceToUnderscore(sportType)))
	}

	for _, v := range event.GameLevels {
		if ruGameLevel, ok := models.EnToRuGameLevel(v); ok {
			hashtags = append(hashtags, fmt.Sprintf(formatHashtag, replaceDashAndSpaceToUnderscore(ruGameLevel)))
		}
	}

	if event.IsFree {
		hashtags = append(hashtags, "#бесплатно")
	} else {
		hashtags = append(hashtags, "#платно")
	}

	return hashtags
}
