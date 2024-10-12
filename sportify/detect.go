package main

import (
	"fmt"
	"regexp"
)

var (
	errInvalidMinMatches = fmt.Errorf("Invalid minMatches value")
	errCompilingError    = fmt.Errorf("Error compiling regular expression")
	errMatchingError     = fmt.Errorf("Error matching regular expression")

	detectRegExps = []string{
		"\b(турнир|соревнование|матч|встреча)\b",
		"\b(игра|поединок|матч)\b.*\b(по|против|с)\b",
		"\b(игра|поединок|матч)\b.*\b(вид спорта|видом спорта|тип спорта|типом спорта)\b",
		"\b(проходит|состоится|начнется|закончится)\b.*\b(в|на)\b.*\b(место|локация|адрес)\b",
		"\b(дата|время|когда)\b.*\b(проходит|состоится|начнется|закончится)\b",
		"\b(требуется|нужен|желателен)\b.*\b(уровень|уровень игрока|уровень игроков)\b",
		"\b(цена|стоимость)\b.*\b(платная|бесплатная|стоит)\b",
		"\b(команда|командир|игрок)\b.*\b(требуется|нужен|желателен)\b",
		"\b(присоединиться|участвовать|участвуйте)\b",
		"\b(тур|матч|игра)\b.*\b(платно|бесплатно)\b",
		"\b(событие|мероприятие)\b.*\b(спортивное|вид спорта|тип спорта)\b",
	}
)

// detect checks if a string satisfies a specified number of regular expressions.
func detect(text string, regexps []string, minMatches int) (bool, error) {
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
		if re.MatchString(text) {
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
