package evaluation_test

import (
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoreRunner_Score(t *testing.T) {
	runner := evaluation.NewScoreRunner(FakeScorer{})
	responses := make(chan evaluation.ChatResponse)
	results := runner.Score(responses)

	response := evaluation.ChatResponse{Query: "How are you?", Response: "Good", Source: "https://me.example.com"}
	responses <- response

	result := <-results
	assert.Equal(t, evaluation.ResponseScore{
		Relevance:       40,
		Correctness:     50,
		AppropriateTone: 60,
		Politeness:      70,
	}, result.Score)
	assert.Equal(t, response, result.Response)
}

type FakeScorer struct {
}

func (f FakeScorer) Score(response evaluation.ChatResponse) (evaluation.ResponseScore, error) {
	return evaluation.ResponseScore{
		Relevance:       40,
		Correctness:     50,
		AppropriateTone: 60,
		Politeness:      70,
	}, nil
}