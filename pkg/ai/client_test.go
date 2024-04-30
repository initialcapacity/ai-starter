package ai_test

import (
	"context"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"log"
	"net/http"
	"testing"
)

func testAiClient(openAiEndpoint string) ai.Client {
	openAiClient, err := azopenai.NewClientForOpenAI(openAiEndpoint, nil, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("unable to create Open AI client: %w", err))
	}

	return ai.Client{OpenAiClient: openAiClient}
}

func TestClient_CreateEmbedding(t *testing.T) {
	endpoint, server := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(http.StatusOK) })
		mux.HandleFunc("/embeddings", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, err := w.Write([]byte(`{
				"data": [
					{ "embedding": [1, 2, 3, 4] }
				]
			}`))
			assert.NoError(t, err)
		})
	})
	defer testsupport.StopTestServer(t, server)

	client := testAiClient(endpoint)
	result, err := client.CreateEmbedding(context.Background(), "some text")

	assert.NoError(t, err)
	assert.Equal(t, []float32{1, 2, 3, 4}, result)
}
