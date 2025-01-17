package scores_test

import (
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoresGateway_Save(t *testing.T) {
	testDb := testsupport.NewTestDb(t)

	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response,
                    		chat_model, temperature, embeddings_model)
						values ('07e87e22-7b55-4023-8c67-64204d30a900', 'You are a bot', 'Hi', 'https://example.com', 'Hello',
							'gpt-11-max', 0.5, 'text-embedding-medium')`)

	scoresGateway := scores.NewGateway(testDb.DB)

	id, err := scoresGateway.Save("07e87e22-7b55-4023-8c67-64204d30a900", 11, 12, 13, 14)
	assert.NoError(t, err)

	saved := testDb.QueryOneMap(`
		select query_response_id,
			(score -> 'relevance')::int as relevance,
			(score -> 'correctness')::int as correctness,
			(score -> 'appropriate_tone')::int as appropriate_tone,
			(score -> 'politeness')::int as politeness,
			score_version
		from response_scores where id = $1`, id)
	assert.NoError(t, err)
	assert.Equal(t, map[string]any{
		"query_response_id": "07e87e22-7b55-4023-8c67-64204d30a900",
		"score_version":     int64(1),
		"relevance":         int64(11),
		"correctness":       int64(12),
		"appropriate_tone":  int64(13),
		"politeness":        int64(14),
	}, saved)
}
