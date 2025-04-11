package testsupport

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/pkg/slicesupport"
	"github.com/pgvector/pgvector-go"
	"strings"
)

func CreateVector(oneIndex int) []float64 {
	embedding := make([]float64, 3072)
	embedding[oneIndex] = 1
	return embedding
}

func CreatePgVector(oneIndex int) pgvector.Vector {
	return pgvector.NewVector(slicesupport.Map(CreateVector(oneIndex), func(i float64) float32 { return float32(i) }))
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
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
