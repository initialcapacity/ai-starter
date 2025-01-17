package evaluation_test

import (
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/pgvector/pgvector-go"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestCannedResponseEvaluator_Run(t *testing.T) {
	testDb := testsupport.NewTestDb(t)

	testDb.Execute("insert into data (id, source, content) values ('aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16', 'https://example.com', 'some content')")
	testDb.Execute("insert into chunks (id, data_id, content) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', 'aaaaaaaa-2f3f-4bc9-8dba-ba397156cc16','a chunk')")
	testDb.Execute("insert into embeddings (chunk_id, embedding) values ('bbbbbbbb-2f3f-4bc9-8dba-ba397156cc16', $1)", pgvector.NewVector(testsupport.CreateVector(0)))
	queryService := query.NewService(analysis.NewEmbeddingsGateway(testDb.DB), testsupport.FakeAi{}, query.NewResponsesGateway(testDb.DB))
	retriever := query.NewChatResponseRetriever(queryService)

	evaluator := evaluation.NewCannedResponseEvaluator(retriever, scores.NewRunner(FakeScorer{}),
		evaluation.NewCSVReporter(), evaluation.NewMarkdownReporter())

	testDirectory := t.TempDir()
	err := evaluator.Run(testDirectory, []string{"Sound good?"})
	assert.NoError(t, err)

	csvContent, err := os.ReadFile(path.Join(testDirectory, "scores.csv"))
	assert.NoError(t, err)
	assert.Equal(t, "Query,Response,Source,Relevance,Correctness,Appropriate Tone,Politeness\nSound good?,Sounds good,https://example.com,40,50,60,70\n", string(csvContent))

	mdContent, err := os.ReadFile(path.Join(testDirectory, "scores.md"))
	assert.NoError(t, err)
	assert.Contains(t, string(mdContent), "Sounds good")
}

type FakeScorer struct {
}

func (f FakeScorer) Score(_ query.ChatResponse) (scores.ResponseScore, error) {
	return scores.ResponseScore{
		Relevance:       40,
		Correctness:     50,
		AppropriateTone: 60,
		Politeness:      70,
	}, nil
}
