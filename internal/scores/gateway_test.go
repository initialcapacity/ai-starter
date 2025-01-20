package scores_test

import (
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/initialcapacity/ai-starter/pkg/slicesupport"
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

func TestScoresGateway_FindForResponseId(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response,
                    		chat_model, temperature, embeddings_model)
						values ('11111111-7b55-4023-8c67-64204d30a900', 'You are a bot', 'Hi', 'https://example.com', 'Hello',
							'gpt-11-max', 0.5, 'text-embedding-medium')`)
	testDb.Execute(`insert into response_scores (id, query_response_id, score, score_version)
						values ('aaaaaaaa-7b55-4023-8c67-64204d30a900', '11111111-7b55-4023-8c67-64204d30a900', '{"relevance": 11, "correctness": 12, "appropriate_tone": 13, "politeness": 14}', 1)`)

	scoresGateway := scores.NewGateway(testDb.DB)

	record, err := scoresGateway.FindForResponseId("11111111-7b55-4023-8c67-64204d30a900")

	assert.NoError(t, err)
	assert.Equal(t, "aaaaaaaa-7b55-4023-8c67-64204d30a900", record.Id)
	assert.Equal(t, "11111111-7b55-4023-8c67-64204d30a900", record.QueryResponseId)
	assert.Equal(t, 11, record.Relevance)
	assert.Equal(t, 12, record.Correctness)
	assert.Equal(t, 13, record.AppropriateTone)
	assert.Equal(t, 14, record.Politeness)
}

func TestScoresGateway_FindForResponseId_NotFound(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	scoresGateway := scores.NewGateway(testDb.DB)

	_, err := scoresGateway.FindForResponseId("bbaaaadd-7b55-4023-8c67-64204d30a900")

	assert.Error(t, err)
}

func TestScoresGateway_ListForResponseIds(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response,
                    		chat_model, temperature, embeddings_model)
						values ('11111111-7b55-4023-8c67-64204d30a900', 'You are a bot', 'Hi', 'https://example.com', 'Hello',
							'gpt-11-max', 0.5, 'text-embedding-medium')`)
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response,
                    		chat_model, temperature, embeddings_model)
						values ('22222222-7b55-4023-8c67-64204d30a900', 'You are a bot', 'Hi', 'https://example.com', 'Hello',
							'gpt-11-max', 0.5, 'text-embedding-medium')`)
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response,
                    		chat_model, temperature, embeddings_model)
						values ('33333333-7b55-4023-8c67-64204d30a900', 'You are a bot', 'Hi', 'https://example.com', 'Hello',
							'gpt-11-max', 0.5, 'text-embedding-medium')`)
	testDb.Execute(`insert into response_scores (id, query_response_id, score, score_version)
						values ('aaaaaaaa-7b55-4023-8c67-64204d30a900', '11111111-7b55-4023-8c67-64204d30a900', '{"relevance": 11, "correctness": 12, "appropriate_tone": 13, "politeness": 14}', 1)`)
	testDb.Execute(`insert into response_scores (id, query_response_id, score, score_version)
						values ('bbbbbbbb-7b55-4023-8c67-64204d30a900', '11111111-7b55-4023-8c67-64204d30a900', '{"relevance": 21, "correctness": 22, "appropriate_tone": 23, "politeness": 24}', 1)`)
	testDb.Execute(`insert into response_scores (id, query_response_id, score, score_version)
						values ('cccccccc-7b55-4023-8c67-64204d30a900', '22222222-7b55-4023-8c67-64204d30a900', '{"relevance": 31, "correctness": 32, "appropriate_tone": 33, "politeness": 34}', 1)`)

	scoresGateway := scores.NewGateway(testDb.DB)

	records, err := scoresGateway.ListForResponseIds([]string{"11111111-7b55-4023-8c67-64204d30a900", "22222222-7b55-4023-8c67-64204d30a900", "33333333-7b55-4023-8c67-64204d30a900"})

	assert.NoError(t, err)
	assert.Len(t, records, 3)
	assert.Equal(t, []string{"aaaaaaaa-7b55-4023-8c67-64204d30a900", "bbbbbbbb-7b55-4023-8c67-64204d30a900", "cccccccc-7b55-4023-8c67-64204d30a900"}, slicesupport.Map(records, func(r scores.ScoreRecord) string { return r.Id }))
	assert.Equal(t, "11111111-7b55-4023-8c67-64204d30a900", records[0].QueryResponseId)
	assert.Equal(t, 11, records[0].Relevance)
	assert.Equal(t, 12, records[0].Correctness)
	assert.Equal(t, 13, records[0].AppropriateTone)
	assert.Equal(t, 14, records[0].Politeness)
}
