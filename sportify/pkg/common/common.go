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
