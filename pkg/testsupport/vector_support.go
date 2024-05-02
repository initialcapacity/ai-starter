package testsupport

func CreateVector(oneIndex int) []float32 {
	embedding := make([]float32, 3072)
	embedding[oneIndex] = 1
	return embedding
}
