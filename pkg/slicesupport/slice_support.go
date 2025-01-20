package slicesupport

func Map[T any, U any](input []T, transformation func(T) U) []U {
	output := make([]U, len(input))
	for i, item := range input {
		output[i] = transformation(item)
	}
	return output
}

func Find[T any](input []T, predicate func(T) bool) (T, bool) {
	for _, item := range input {
		if predicate(item) {
			return item, true
		}
	}
	var emptyResult T
	return emptyResult, false
}
