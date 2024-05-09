package testsupport

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
)

func CreateVector(oneIndex int) []float32 {
	embedding := make([]float32, 3072)
	embedding[oneIndex] = 1
	return embedding
}

type Number interface {
	constraints.Integer | constraints.Float
}

func VectorToString[T Number](vector []T) string {
	builder := strings.Builder{}

	builder.WriteString("[")
	for i, v := range vector {
		builder.WriteString(fmt.Sprint(v))
		if i < len(vector)-1 {
			builder.WriteString(", ")
		}
	}
	builder.WriteString("]")

	return builder.String()
}
