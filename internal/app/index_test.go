package app_test

import (
	"github.com/initialcapacity/ai-starter/internal/app"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestIndex_Get(t *testing.T) {
	appEndpoint := testsupport.StartTestServer(t, app.Handlers(testsupport.NewTestAiClient(""), nil))

	resp, err := http.Get(appEndpoint)
	assert.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := string(bytes)
	assert.Contains(t, body, "AI Starter")
}

func TestIndex_Post(t *testing.T) {
	aiEndpoint := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.HandleGetStreamCompletion(mux, "Sounds good")
		testsupport.HandleCreateEmbedding(mux, testsupport.CreateVector(0))
	})
	testDb := testsupport.NewTestDb(t)

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")
	testDb.Execute("insert into embeddings (chunk_id, embedding) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', $1)", testsupport.CreatePgVector(0))

	appEndpoint := testsupport.StartTestServer(t, app.Handlers(testsupport.NewTestAiClient(aiEndpoint), testDb.DB))

	resp, err := http.Post(appEndpoint, "application/x-www-form-urlencoded", strings.NewReader("query=what%20do%20you%20think"))
	assert.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := string(bytes)
	assert.Contains(t, body, "what do you think")
	assert.Contains(t, body, "https://example.com")
	assert.Contains(t, body, "Sounds good")
}
