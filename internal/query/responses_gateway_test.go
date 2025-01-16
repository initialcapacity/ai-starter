package query_test

import (
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestResponsesGateway_Create(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	gateway := query.NewResponsesGateway(testDb.DB)

	record, err := gateway.Create("hello", "what's up?", "https://example.com", "Not much", "gpt-42", "text-embeddings-test", 1)
	require.NoError(t, err)

	assert.Equal(t, "hello", record.SystemPrompt)
	assert.Equal(t, "what's up?", record.UserQuery)
	assert.Equal(t, "https://example.com", record.Source)
	assert.Equal(t, "Not much", record.Response)
	assert.Equal(t, "gpt-42", record.ChatModel)
	assert.Equal(t, "text-embeddings-test", record.EmbeddingsModel)
	assert.Equal(t, float32(1), record.Temperature)

	result := testDb.QueryOneMap("select system_prompt, user_query, source, response, chat_model, embeddings_model, temperature from query_responses where id = $1", record.Id)
	assert.Equal(t, map[string]any{
		"system_prompt":    "hello",
		"user_query":       "what's up?",
		"source":           "https://example.com",
		"response":         "Not much",
		"chat_model":       "gpt-42",
		"embeddings_model": "text-embeddings-test",
		"temperature":      float64(1),
	}, result)
}

func TestResponsesGateway_List(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	gateway := query.NewResponsesGateway(testDb.DB)

	testDb.Execute("insert into query_responses (id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature) values ('11111111-3c62-4174-a53c-f317f49ba2fa', 'hi', 'what is going on?', 'https://example.com/1', 'A lot', 'gpt-43', 'text-embeddings-test', 1.5)")
	testDb.Execute("insert into query_responses (id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature) values ('22222222-3c62-4174-a53c-f317f49ba2fa', 'hello', 'what is up?', 'https://example.com/2', 'Not much', 'gpt-42', 'text-embeddings-test', 1)")

	records, err := gateway.List()
	require.NoError(t, err)

	assert.Equal(t, 2, len(records))
	assert.Equal(t, records[0].Id, "22222222-3c62-4174-a53c-f317f49ba2fa")
	assert.Equal(t, "hello", records[0].SystemPrompt)
	assert.Equal(t, "what is up?", records[0].UserQuery)
	assert.Equal(t, "https://example.com/2", records[0].Source)
	assert.Equal(t, "Not much", records[0].Response)
	assert.Equal(t, "gpt-42", records[0].ChatModel)
	assert.Equal(t, "text-embeddings-test", records[0].EmbeddingsModel)
	assert.Equal(t, float32(1), records[0].Temperature)
	assert.Equal(t, records[1].Id, "11111111-3c62-4174-a53c-f317f49ba2fa")
}

func TestResponsesGateway_Find(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	gateway := query.NewResponsesGateway(testDb.DB)

	testDb.Execute("insert into query_responses (id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature) values ('11111111-3c62-4174-a53c-f317f49ba2fa', 'hello', 'what is up?', 'https://example.com', 'Not much', 'gpt-42', 'text-embeddings-test', 1)")

	record, err := gateway.Find("11111111-3c62-4174-a53c-f317f49ba2fa")
	require.NoError(t, err)

	assert.Equal(t, record.Id, "11111111-3c62-4174-a53c-f317f49ba2fa")
	assert.Equal(t, "hello", record.SystemPrompt)
	assert.Equal(t, "what is up?", record.UserQuery)
	assert.Equal(t, "https://example.com", record.Source)
	assert.Equal(t, "Not much", record.Response)
	assert.Equal(t, "gpt-42", record.ChatModel)
	assert.Equal(t, "text-embeddings-test", record.EmbeddingsModel)
	assert.Equal(t, float32(1), record.Temperature)
}

func TestResponsesGateway_FindNotFound(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	gateway := query.NewResponsesGateway(testDb.DB)

	_, err := gateway.Find("bbaaaadd-3c62-4174-a53c-f317f49ba2fa")
	require.Error(t, err)
}

func TestResponsesGateway_ListMissingScores(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	gateway := query.NewResponsesGateway(testDb.DB)

	testDb.Execute("insert into query_responses (id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature) values ('11111111-3c62-4174-a53c-f317f49ba2fa', 'hi', 'what is going on?', 'https://example.com/1', 'A lot', 'gpt-43', 'text-embeddings-test', 1.5)")
	testDb.Execute("insert into query_responses (id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature) values ('22222222-3c62-4174-a53c-f317f49ba2fa', 'hello', 'what is up?', 'https://example.com/2', 'Not much', 'gpt-42', 'text-embeddings-test', 1)")
	testDb.Execute("insert into response_scores (query_response_id, score, score_version) values ('11111111-3c62-4174-a53c-f317f49ba2fa', '{}', 11)")

	records, err := gateway.ListMissingScores()
	require.NoError(t, err)

	assert.Equal(t, 1, len(records))
	assert.Equal(t, records[0].Id, "22222222-3c62-4174-a53c-f317f49ba2fa")
	assert.Equal(t, "hello", records[0].SystemPrompt)
	assert.Equal(t, "what is up?", records[0].UserQuery)
	assert.Equal(t, "https://example.com/2", records[0].Source)
	assert.Equal(t, "Not much", records[0].Response)
	assert.Equal(t, "gpt-42", records[0].ChatModel)
	assert.Equal(t, "text-embeddings-test", records[0].EmbeddingsModel)
	assert.Equal(t, float32(1), records[0].Temperature)
}
