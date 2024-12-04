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

func HandleGetStreamCompletion(mux *http.ServeMux, response string) {
	mux.HandleFunc("POST /chat/completions", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(Stream(fmt.Sprintf(`{ "choices": [ { "delta": { "role": "assistant", "content": "%s" } } ] }`, response)))
	})
}

func HandleGetCompletion(mux *http.ServeMux, response string) {
	Handle(mux, "POST /chat/completions", fmt.Sprintf(`{ "choices": [ { "message": { "role": "assistant", "content": "%s" } } ] }`, response))
}
