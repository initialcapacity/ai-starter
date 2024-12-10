package app_test

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/app"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestQueryResponses_Get(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response, model, temperature)
			values ('11111111-2f3f-4bc9-8dba-ba397156cc16', 'Hello', 'Say hi', 'https://example.com', 'Hi there', 'gpt-11', 1.2)`)

	server := websupport.NewServer(app.Handlers(testsupport.NewTestAiClient(""), testDb.DB))
	port, _ := server.Start("localhost", 0)
	defer func(server *websupport.Server) {
		_ = server.Stop()
	}(server)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/query_responses", port))
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
	assert.Contains(t, body, "1.2")
}

func TestShowQueryResponse_Get(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	testDb.Execute(`insert into query_responses (id, system_prompt, user_query, source, response, model, temperature)
			values ('11111111-2f3f-4bc9-8dba-ba397156cc16', 'Hello', 'Say hi', 'https://example.com', 'Hi there', 'gpt-11', 1.2)`)

	server := websupport.NewServer(app.Handlers(testsupport.NewTestAiClient(""), testDb.DB))
	port, _ := server.Start("localhost", 0)
	defer func(server *websupport.Server) {
		_ = server.Stop()
	}(server)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/query_responses/11111111-2f3f-4bc9-8dba-ba397156cc16", port))
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
	assert.Contains(t, body, "1.2")
}
