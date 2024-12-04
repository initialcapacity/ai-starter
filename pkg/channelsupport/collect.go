package channelsupport

func CollectSlice[T any](items chan T) []T {
	result := make([]T, 0)

	for item := range items {
		result = append(result, item)
	}

	return result
}
