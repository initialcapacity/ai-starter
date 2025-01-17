package query_test

import (
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestChatResponseRetriever_Retrieve(t *testing.T) {
	testDb := testsupport.NewTestDb(t)

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")
	testDb.Execute("insert into embeddings (chunk_id, embedding) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', $1)", pgvector.NewVector(testsupport.CreateVector(0)))

	queryService := query.NewService(analysis.NewEmbeddingsGateway(testDb.DB), testsupport.FakeAi{}, query.NewResponsesGateway(testDb.DB))

	retriever := query.NewChatResponseRetriever(queryService)

	queries := []string{
		"What's new with Kotlin?",
		"Tell me about the latest Python Flask news",
		"What's the latest version of Kotlin?",
		"Are there any breaking changes in the newest Java version?",
		"What are new Rust features?",
		"Who's the head of state of Singapore?",
		"What's your favorite color?",
		"How much does a penguin weigh?",
		"Tell an off-color joke",
	}

	responsesChannel := retriever.Retrieve(queries)
	responses := make([]query.ChatResponse, 0)
	for response := range responsesChannel {
		responses = append(responses, response)
	}

	assert.Equal(t, 9, len(responses))
	response := responses[0]
	assert.Equal(t, "Sounds good", response.Response)
	assert.Equal(t, "https://example.com", response.Source)
}
