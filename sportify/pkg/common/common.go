package common

func Ref[T any](val T) *T {
	return &val
}

func Find[T any](collection []T, predicate func(item T) bool) (T, bool) {
	for i := range collection {
		if predicate(collection[i]) {
			return collection[i], true
		}
	}

	return *new(T), false
}

func Map[T any, V any](mapping func(T) V, elements []T) []V {
	result := make([]V, 0, len(elements))
	for _, arg := range elements {
		result = append(result, mapping(arg))
	}
	return result
}

func NewValWithFallback[T any](newVal, fallbackVal *T) T {
	if newVal != nil {
		return *newVal
	}
	return *fallbackVal
}
