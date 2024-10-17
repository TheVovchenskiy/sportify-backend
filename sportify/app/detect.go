package app

import (
	"fmt"
	"regexp"
)

var (
	errInvalidMinMatches = fmt.Errorf("Invalid minMatches value")
	errCompilingError    = fmt.Errorf("Error compiling regular expression")
	errMatchingError     = fmt.Errorf("Error matching regular expression")

	// FIXME: fix possible regex injections
	SportEventRegExps = []string{
		`(?:манеж|стадион|метро|парк|поле).*?(?:«[^»]+»|".+?"|[А-ЯЁа-яё]+)`,                                                           // place
		`(?:\d{1,2}\s*(?:января|февраля|марта|апреля|мая|июня|июля|августа|сентября|октября|ноября|декабря)|\d{1,2}.\d{1,2}.\d{2,4})`, // date
		`\d{1,2}:\d{2}(?:-\d{1,2}:\d{2})?`,                                       // time
		`(?:\d+\s*команд[ыа]?\s*\d+(×|\*|на|x|х)\d+|\d+(×|\*|на|x|х)\d+)`,        // game format
		`\d+(?:.\d+)?\s*(?:час[аов]|минут[ыа])`,                                  // duration
		`(взнос|цена)\s*\d+\s*(?:р.?|рубл(?:ей|я)?)|\d+\s*(?:р.?|рубл(?:ей|я)?)`, // price
		`(?:требуется|нужно|не хватает)\s*\d+\s*(?:человек|игрок[ов]|мест[ао])`,  // level
		`(?:видео(?:съёмка)?|вода для игроков|душевые|раздевалки)`,               // facilities
		`(?:играем|поиграть|матч|тренировка)`,                                    // type
		`https?:\/\/[^\s]+`, // link

		"(турнир|соревнование|матч|встреча)",
		"(игра|поединок|матч).*(по|против|с)",
		"(игра|поединок|матч).*(вид спорта|видом спорта|тип спорта|типом спорта)",
		"(проходит|состоится|начнется|закончится).*(в|на).*(место|локация|адрес)",
		"(дата|время|когда).*(проходит|состоится|начнется|закончится)",
		"(требуется|нужен|желателен).*(уровень|уровень игрока|уровень игроков)",
		"(цена|стоимость|взнос|плата)",
		"(платная|бесплатная|стоит)",
		"(команда|командир|игрок).*(требуется|нужен|желателен)",
		"(присоединиться|участвовать|участвуйте)",
		"(тур|матч|игра).*(платно|бесплатно)",
		"(событие|мероприятие).*(спортивное|вид спорта|тип спорта)",
	}
)

// DetectEventMessage checks if a string satisfies a specified number of regular expressions.
func (a *App) DetectEventMessage(text string, regexps []string, minMatches int) (bool, error) {
	if minMatches < 0 {
		return false, fmt.Errorf("%w: must be >= 0", errInvalidMinMatches)
	}
	if minMatches == 0 {
		return true, nil
	}

	matched := 0

	for _, regex := range regexps {
		re, err := regexp.Compile(regex)
		if err != nil {
			return false, fmt.Errorf("%w %q: %w", errCompilingError, regex, err)
		}
		if re.MatchString("(?i)" + text) {
			matched++
		}
		if matched >= minMatches {
			return true, nil
		}
	}

	if matched < minMatches {
		return false, fmt.Errorf(
			"%w: the string does not satisfy the specified minimum number of regular expressions (%d of %d)",
			errMatchingError,
			matched,
			minMatches,
		)
	}

	return true, nil
}
