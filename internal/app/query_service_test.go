package app_test

import (
	"context"
	"errors"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/internal/app"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryService_FetchResponse(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	insertData(testDb)

	service := app.NewQueryService(analyzer.NewEmbeddingsGateway(testDb.DB), fakeAi{})

	result, err := service.FetchResponse(context.Background(), "Does this sound good?")
	assert.NoError(t, err)
	message := <-result.Response

	assert.Equal(t, "https://example.com", result.Source)
	assert.Equal(t, "Sounds good", message)
}

func TestQueryService_FetchResponse_EmbeddingError(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	insertData(testDb)
	service := app.NewQueryService(analyzer.NewEmbeddingsGateway(testDb.DB), fakeAi{embeddingError: errors.New("bad news")})

	_, err := service.FetchResponse(context.Background(), "Does this sound good?")

	assert.EqualError(t, err, "bad news")
}

func TestQueryService_FetchResponse_NoEmbeddings(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	service := app.NewQueryService(analyzer.NewEmbeddingsGateway(testDb.DB), fakeAi{})

	_, err := service.FetchResponse(context.Background(), "Does this sound good?")

	assert.EqualError(t, err, "sql: no rows in result set")
}

func TestQueryService_FetchResponse_CompletionError(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	defer testDb.Close()
	insertData(testDb)
	service := app.NewQueryService(analyzer.NewEmbeddingsGateway(testDb.DB), fakeAi{completionError: errors.New("bad news")})

	_, err := service.FetchResponse(context.Background(), "Does this sound good?")

	assert.EqualError(t, err, "bad news")
}

func insertData(testDb *testsupport.TestDb) {
	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")
	testDb.Execute("insert into embeddings (chunk_id, embedding) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', $1)", pgvector.NewVector(testsupport.CreateVector(0)))

}

type fakeAi struct {
	embeddingError  error
	completionError error
}

func (f fakeAi) CreateEmbedding(_ context.Context, _ string) ([]float32, error) {
	if f.completionError != nil {
		return nil, f.completionError
	}

	return testsupport.CreateVector(0), nil
}

func (f fakeAi) GetChatCompletion(_ context.Context, _ []ai.ChatMessage) (chan string, error) {
	if f.embeddingError != nil {
		return nil, f.embeddingError
	}

	response := make(chan string)
	go func() {
		response <- "Sounds good"
	}()
	return response, nil
}
