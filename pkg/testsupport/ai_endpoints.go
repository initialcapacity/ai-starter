package testsupport

import (
	"fmt"
	"net/http"
)

func HandleCreateEmbedding[T Number](mux *http.ServeMux, vector []T) {
	Handle(mux, "POST /embeddings", fmt.Sprintf(`{
			"data": [
				{ "embedding": %s }
			]
		}`, VectorToString(vector)))
}
