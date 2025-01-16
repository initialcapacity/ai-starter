package evaluation_test

import (
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPastResponseEvaluator_Run(t *testing.T) {
	testDb := testsupport.NewTestDb(t)

	responsesGateway := query.NewResponsesGateway(testDb.DB)
	scoresGateway := evaluation.NewScoresGateway(testDb.DB)
	evaluator := evaluation.NewPastResponseEvaluator(responsesGateway, scoresGateway, FakeScorer{})

	response, err := responsesGateway.Create("You are a bot", "Hi", "https://example.com", "Hello", "gpt11", "embeddings-xl", .4)
	assert.NoError(t, err)

	err = evaluator.Run()
	assert.NoError(t, err)

	scores := testDb.QueryMap(`
		select query_response_id,
			(score -> 'relevance')::int as relevance,
			(score -> 'correctness')::int as correctness,
			(score -> 'appropriate_tone')::int as appropriate_tone,
			(score -> 'politeness')::int as politeness,
			score_version
		from response_scores`)

	assert.Len(t, scores, 1)
	assert.Equal(t, map[string]any{
		"query_response_id": response.Id,
		"score_version":     int64(1),
		"relevance":         int64(40),
		"correctness":       int64(50),
		"appropriate_tone":  int64(60),
		"politeness":        int64(70),
	}, scores[0])
}
