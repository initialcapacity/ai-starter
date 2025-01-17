package query_test

import (
	"context"
	"errors"
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestQueryService_FetchResponse(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	insertData(testDb)

	responsesGateway := query.NewResponsesGateway(testDb.DB)
	service := query.NewService(analysis.NewEmbeddingsGateway(testDb.DB), testsupport.FakeAi{}, responsesGateway)

	result, err := service.FetchResponse(context.Background(), "Does this sound good?")
	assert.NoError(t, err)

	var parts strings.Builder
	for part := range result.Response {
		parts.WriteString(part)
	}
	message := parts.String()

	assert.Equal(t, "https://example.com", result.Source)
	assert.Equal(t, "Sounds good", message)

	responses, err := responsesGateway.List()
	assert.NoError(t, err)
	assert.Equal(t, 1, len(responses))
	assert.Equal(t, `You are a reporter for a major world newspaper.
Write your response as if you were writing a short, high-quality news article for your paper. Limit your response to one paragraph.

Use the following article for context: a chunk`, responses[0].SystemPrompt)
	assert.Equal(t, "Does this sound good?", responses[0].UserQuery)
	assert.Equal(t, "https://example.com", responses[0].Source)
	assert.Equal(t, "Sounds good", responses[0].Response)
	assert.Equal(t, "gpt-123", responses[0].ChatModel)
	assert.Equal(t, float32(2), responses[0].Temperature)
}

func TestQueryService_FetchResponse_EmbeddingError(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	insertData(testDb)
	service := query.NewService(analysis.NewEmbeddingsGateway(testDb.DB), testsupport.FakeAi{EmbeddingError: errors.New("bad news")}, query.NewResponsesGateway(testDb.DB))

	_, err := service.FetchResponse(context.Background(), "Does this sound good?")

	assert.EqualError(t, err, "bad news")
}

func TestQueryService_FetchResponse_NoEmbeddings(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	service := query.NewService(analysis.NewEmbeddingsGateway(testDb.DB), testsupport.FakeAi{}, query.NewResponsesGateway(testDb.DB))

	_, err := service.FetchResponse(context.Background(), "Does this sound good?")

	assert.EqualError(t, err, "sql: no rows in result set")
}

func TestQueryService_FetchResponse_CompletionError(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	insertData(testDb)
	service := query.NewService(analysis.NewEmbeddingsGateway(testDb.DB), testsupport.FakeAi{CompletionError: errors.New("bad news")}, query.NewResponsesGateway(testDb.DB))

	_, err := service.FetchResponse(context.Background(), "Does this sound good?")

	assert.EqualError(t, err, "bad news")
}

func insertData(testDb *testsupport.TestDb) {
	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")
	testDb.Execute("insert into embeddings (chunk_id, embedding) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', $1)", pgvector.NewVector(testsupport.CreateVector(0)))
}
