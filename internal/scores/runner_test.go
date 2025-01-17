package scores_test

import (
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoreRunner_Score(t *testing.T) {
	runner := scores.NewRunner(FakeScorer{})
	response := query.ChatResponse{Query: "How are you?", Response: "Good", Source: "https://me.example.com"}

	responses := make(chan query.ChatResponse, 1)
	responses <- response
	close(responses)

	results := runner.Score(responses)

	assert.Len(t, results, 1)
	result := results[0]
	assert.Equal(t, scores.ResponseScore{
		Relevance:       40,
		Correctness:     50,
		AppropriateTone: 60,
		Politeness:      70,
	}, result.Score)
	assert.Equal(t, response, result.Response)
}

type FakeScorer struct {
}

func (f FakeScorer) Score(_ query.ChatResponse) (scores.ResponseScore, error) {
	return scores.ResponseScore{
		Relevance:       40,
		Correctness:     50,
		AppropriateTone: 60,
		Politeness:      70,
	}, nil
}
