package slicesupport

func Map[T any, U any](input []T, transformation func(T) U) []U {
	output := make([]U, len(input))
	for i, item := range input {
		output[i] = transformation(item)
	}
	return output
}
