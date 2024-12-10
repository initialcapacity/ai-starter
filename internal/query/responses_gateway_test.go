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
	defer testDb.Close()
	gateway := query.NewResponsesGateway(testDb.DB)

	record, err := gateway.Create("hello", "what's up?", "https://example.com", "Not much", "gpt-42", 1)
	require.NoError(t, err)

	assert.Equal(t, "hello", record.SystemPrompt)
	assert.Equal(t, "what's up?", record.UserQuery)
	assert.Equal(t, "https://example.com", record.Source)
	assert.Equal(t, "Not much", record.Response)
	assert.Equal(t, "gpt-42", record.Model)
	assert.Equal(t, float32(1), record.Temperature)

	result := testDb.QueryOneMap("select system_prompt, user_query, source, response, model, temperature from query_responses where id = $1", record.Id)
	assert.Equal(t, map[string]any{
		"system_prompt": "hello",
		"user_query":    "what's up?",
		"source":        "https://example.com",
		"response":      "Not much",
		"model":         "gpt-42",
		"temperature":   float64(1),
	}, result)
}

func TestResponsesGateway_List(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	gateway := query.NewResponsesGateway(testDb.DB)

	testDb.Execute("insert into query_responses (id, system_prompt, user_query, source, response, model, temperature) values ('11111111-3c62-4174-a53c-f317f49ba2fa', 'hi', 'what is going on?', 'https://example.com/1', 'A lot', 'gpt-43', 1.5)")
	testDb.Execute("insert into query_responses (id, system_prompt, user_query, source, response, model, temperature) values ('22222222-3c62-4174-a53c-f317f49ba2fa', 'hello', 'what is up?', 'https://example.com/2', 'Not much', 'gpt-42', 1)")

	records, err := gateway.List()
	require.NoError(t, err)

	assert.Equal(t, 2, len(records))
	assert.Equal(t, records[0].Id, "22222222-3c62-4174-a53c-f317f49ba2fa")
	assert.Equal(t, "hello", records[0].SystemPrompt)
	assert.Equal(t, "what is up?", records[0].UserQuery)
	assert.Equal(t, "https://example.com/2", records[0].Source)
	assert.Equal(t, "Not much", records[0].Response)
	assert.Equal(t, "gpt-42", records[0].Model)
	assert.Equal(t, float32(1), records[0].Temperature)
	assert.Equal(t, records[1].Id, "11111111-3c62-4174-a53c-f317f49ba2fa")
}
