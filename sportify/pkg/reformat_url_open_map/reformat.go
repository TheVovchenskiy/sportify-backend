package reformat_url_open_map

import "strings"

func replaceCityURLOpenMap(url string) string {
	separator := "г "
	idxCityStart := strings.Index(url, separator)
	idxCityEnd := idxCityStart + len(separator)

	return url[idxCityEnd:]
}

func replaceHomeURLOpenMap(url string) string {
	url = strings.ReplaceAll(url, " д ", " дом ")

	separator := " дом "
	idxHomeStart := strings.Index(url, separator)
	if idxHomeStart == -1 {
		return url
	}

	idxHomeEnd := idxHomeStart + len(separator)
	endWithHomeWithExtra := url[idxHomeEnd:]

	idxExtraStart := strings.Index(endWithHomeWithExtra, " ")
	if idxExtraStart == -1 {
		url = strings.ReplaceAll(url, "дом ", "")
		return url
	}

	url = url[:idxHomeEnd+idxExtraStart]
	url = strings.ReplaceAll(url, "дом ", "")

	return url
}

func replaceAnotherURLOpenMap(url string) string {
	replacements := map[string]string{
		" ал":     " аллея",
		" б-р":    " бульвар",
		" взв":    " взвоз",
		" взд":    " въезд",
		" дор":    " дорога",
		" ззд":    " заезд",
		" км":     " километр",
		" к-цо":   " кольцо",
		" коса":   " коса",
		" лн":     " линия",
		" мгстр":  " магистраль",
		" наб":    " набережная",
		" пер-д":  " переезд",
		" пер":    " переулок",
		" пл-ка":  " площадка",
		" пл":     " площадь",
		" пр-д":   " проезд",
		" пр-кт":  " проспект",
		" пр-ка":  " просека",
		" пр-к":   " просек",
		" пр-лок": " проселок",
		" проул":  " проулок",
		" рзд":    " разъезд",
		" ряд":    " ряд",
		" с-р":    " сквер",
		" с-к":    " спуск",
		" сзд":    " съезд",
		" тракт":  " тракт",
		" туп":    " тупик",
		" ул":     " улица",
		" ш":      " шоссе",
		" тер":    " территория",
	}

	for abbrev, full := range replacements {
		url = strings.ReplaceAll(url, abbrev, full)
	}

	url = strings.ReplaceAll(url, ".", "")

	return url
}

func ReformatURLOpenMap(url string) string {
	url = replaceCityURLOpenMap(url)
	url = replaceHomeURLOpenMap(url)
	url = replaceAnotherURLOpenMap(url)

	return url
}
