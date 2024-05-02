package ai_test

import (
	"context"
	"fmt"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_CreateEmbedding(t *testing.T) {
	endpoint, server := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.Handle(mux, "/embeddings", `{
				"data": [
					{ "embedding": [1, 2, 3, 4] }
				]
			}`)
	})
	defer testsupport.StopTestServer(t, server)

	client := testsupport.NewTestAiClient(endpoint)
	result, err := client.CreateEmbedding(context.Background(), "some text")

	assert.NoError(t, err)
	assert.Equal(t, []float32{1, 2, 3, 4}, result)
}

func TestClient_GetChatCompletion(t *testing.T) {
	endpoint, server := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		mux.HandleFunc("/chat/completions", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write(data(`{ "choices": [ { "delta": { "role": "assistant", "content": "Sounds good" } } ] }`))
			assert.NoError(t, err)
		})
	})
	defer testsupport.StopTestServer(t, server)

	client := testsupport.NewTestAiClient(endpoint)
	completion, err := client.GetChatCompletion(context.Background(), []ai.ChatMessage{})

	assert.NoError(t, err)
	assert.Equal(t, "Sounds good", <-completion)
}

func data(json string) []byte {
	return []byte(fmt.Sprintf("data:%s\ndata:[DONE]", json))
}
