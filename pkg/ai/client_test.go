package ai_test

import (
	"context"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestClient_CreateEmbedding(t *testing.T) {
	endpoint := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.HandleCreateEmbedding(mux, []int{1, 2, 3, 4})
	})

	client := testsupport.NewTestAiClient(endpoint)
	result, err := client.CreateEmbedding(context.Background(), "some text")

	assert.NoError(t, err)
	assert.Equal(t, []float64{1, 2, 3, 4}, result)
}

func TestClient_GetChatCompletion(t *testing.T) {
	endpoint := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.HandleGetStreamCompletion(mux, "Sounds good")
	})

	client := testsupport.NewTestAiClient(endpoint)
	completion, err := client.GetChatCompletion(context.Background(), []ai.ChatMessage{})

	assert.NoError(t, err)
	assert.Equal(t, "Sounds good", <-completion)
}

type TestResponse struct {
	Some string `json:"some" jsonschema_description:"A response"`
}

func TestClient_GetJsonChatCompletion(t *testing.T) {
	endpoint := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.HandleGetCompletion(mux, `{\"some\": \"json\"}`)
	})

	client := testsupport.NewTestAiClient(endpoint)
	completion, err := ai.GetJsonChatCompletion[TestResponse](context.Background(), client, []ai.ChatMessage{}, "someName", "a description")

	assert.NoError(t, err)
	assert.Equal(t, TestResponse{Some: "json"}, completion)
}
