package evaluation_test

import (
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestAiScorer_Score(t *testing.T) {
	endpoint := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.HandleGetCompletion(mux,
			`{ \"Relevance\": 10, \"Correctness\": 20, \"AppropriateTone\": 30, \"Politeness\": 40 }`,
		)
	})
	client := testsupport.NewTestAiClient(endpoint)
	scorer := evaluation.NewAiScorer(client)

	score, err := scorer.Score(evaluation.ChatResponse{
		Query:    "Why is the sky blue",
		Response: "Because I said so",
		Source:   "https://sky.example.com",
	})

	assert.NoError(t, err)
	assert.Equal(t, evaluation.ResponseScore{
		Relevance:       10,
		Correctness:     20,
		AppropriateTone: 30,
		Politeness:      40,
	}, score)
}
