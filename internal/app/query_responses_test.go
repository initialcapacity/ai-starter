package app_test

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/app"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestQueryResponses_Get(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature)
			values ('11111111-2f3f-4bc9-8dba-ba397156cc16', 'Hello', 'Say hi', 'https://example.com/1', 'Hi there', 'gpt-11','text-embeddings-test-1',  1.2)`)
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature)
			values ('22222222-2f3f-4bc9-8dba-ba397156cc16', 'Bye', 'Say bye', 'https://example.com/2', 'Bye then', 'gpt-12','text-embeddings-test-2',  1.4)`)
	testDb.Execute(`insert into response_scores (query_response_id, score, score_version)
						values ('11111111-2f3f-4bc9-8dba-ba397156cc16', '{"relevance": 31, "correctness": 32, "appropriate_tone": 33, "politeness": 34}', 1)`)

	appEndpoint := testsupport.StartTestServer(t, app.Handlers(testsupport.NewTestAiClient(""), testDb.DB))

	resp, err := http.Get(fmt.Sprintf("%s/query_responses", appEndpoint))
	assert.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := string(bytes)
	assert.Contains(t, body, "11111111-2f3f-4bc9-8dba-ba397156cc16")
	assert.Contains(t, body, "Hello")
	assert.Contains(t, body, "Say hi")
	assert.Contains(t, body, "https://example.com/1")
	assert.Contains(t, body, "Hi there")
	assert.Contains(t, body, "gpt-11")
	assert.Contains(t, body, "text-embeddings-test-1")
	assert.Contains(t, body, "1.2")
	assert.Contains(t, body, "31")
	assert.Contains(t, body, "32")
	assert.Contains(t, body, "33")
	assert.Contains(t, body, "34")
}

func TestShowQueryResponse_Get(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature)
			values ('11111111-2f3f-4bc9-8dba-ba397156cc16', 'Hello', 'Say hi', 'https://example.com', 'Hi there', 'gpt-11', 'text-embeddings-test', 1.2)`)

	appEndpoint := testsupport.StartTestServer(t, app.Handlers(testsupport.NewTestAiClient(""), testDb.DB))

	resp, err := http.Get(fmt.Sprintf("%s/query_responses/11111111-2f3f-4bc9-8dba-ba397156cc16", appEndpoint))
	assert.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := string(bytes)
	assert.Contains(t, body, "11111111-2f3f-4bc9-8dba-ba397156cc16")
	assert.Contains(t, body, "Hello")
	assert.Contains(t, body, "Say hi")
	assert.Contains(t, body, "https://example.com")
	assert.Contains(t, body, "Hi there")
	assert.Contains(t, body, "gpt-11")
	assert.Contains(t, body, "text-embeddings-test")
	assert.Contains(t, body, "1.2")
}

func TestShowQueryResponse_Get_WithScores(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature)
			values ('11111111-2f3f-4bc9-8dba-ba397156cc16', 'Hello', 'Say hi', 'https://example.com', 'Hi there', 'gpt-11', 'text-embeddings-test', 1.2)`)
	testDb.Execute(`insert into response_scores (query_response_id, score, score_version)
						values ('11111111-2f3f-4bc9-8dba-ba397156cc16', '{"relevance": 31, "correctness": 32, "appropriate_tone": 33, "politeness": 34}', 1)`)

	appEndpoint := testsupport.StartTestServer(t, app.Handlers(testsupport.NewTestAiClient(""), testDb.DB))

	resp, err := http.Get(fmt.Sprintf("%s/query_responses/11111111-2f3f-4bc9-8dba-ba397156cc16", appEndpoint))
	assert.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := string(bytes)
	assert.Contains(t, body, "31")
	assert.Contains(t, body, "32")
	assert.Contains(t, body, "33")
	assert.Contains(t, body, "34")
}
