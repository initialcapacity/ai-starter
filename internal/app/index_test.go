package app_test

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/app"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestIndex_Get(t *testing.T) {
	server := websupport.NewServer(app.Handlers(testsupport.NewTestAiClient(""), nil))
	port, _ := server.Start("localhost", 0)
	defer func(server *websupport.Server) {
		_ = server.Stop()
	}(server)

	resp, err := http.Get(fmt.Sprintf("http://localhost:%d", port))
	assert.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := string(bytes)
	assert.Contains(t, body, "AI Starter")
}

func TestIndex_Post(t *testing.T) {
	aiEndpoint, aiServer := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.HandleGetStreamCompletion(mux, "Sounds good")
		testsupport.HandleCreateEmbedding(mux, testsupport.CreateVector(0))
	})
	defer testsupport.StopTestServer(t, aiServer)

	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")
	testDb.Execute("insert into embeddings (chunk_id, embedding) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', $1)", pgvector.NewVector(testsupport.CreateVector(0)))

	server := websupport.NewServer(app.Handlers(testsupport.NewTestAiClient(aiEndpoint), testDb.DB))
	port, _ := server.Start("localhost", 0)
	defer func(server *websupport.Server) {
		_ = server.Stop()
	}(server)

	resp, err := http.Post(fmt.Sprintf("http://localhost:%d", port), "application/x-www-form-urlencoded", strings.NewReader("query=what%20do%20you%20think"))
	assert.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	body := string(bytes)
	assert.Contains(t, body, "what do you think")
	assert.Contains(t, body, "https://example.com")
	assert.Contains(t, body, "Sounds good")
}
