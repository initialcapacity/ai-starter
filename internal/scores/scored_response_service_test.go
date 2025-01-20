package scores_test

import (
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/initialcapacity/ai-starter/pkg/slicesupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestScoredResponsesService_Find(t *testing.T) {
	service, gateway, responsesGateway := createService(t)
	response, err := responsesGateway.Create("You are a bot", "Hi", "https://example.com", "Hello", "gpt-11-max", "text-embedding-medium", 0.5)
	require.NoError(t, err)
	_, err = gateway.Save(response.Id, 11, 12, 13, 14)
	require.NoError(t, err)

	result, found := service.Find(response.Id)

	assert.True(t, found)
	assert.Equal(t, response.Id, result.Id)
	assert.Equal(t, "You are a bot", result.SystemPrompt)
	assert.Equal(t, "Hi", result.UserQuery)
	assert.Equal(t, "https://example.com", result.Source)
	assert.Equal(t, "Hello", result.Response)
	assert.Equal(t, "gpt-11-max", result.ChatModel)
	assert.Equal(t, "text-embedding-medium", result.EmbeddingsModel)
	assert.Equal(t, float32(0.5), result.Temperature)
	assert.Equal(t, 11, result.Score.Relevance)
	assert.Equal(t, 12, result.Score.Correctness)
	assert.Equal(t, 13, result.Score.AppropriateTone)
	assert.Equal(t, 14, result.Score.Politeness)
}

func TestScoredResponsesService_Find_NoScore(t *testing.T) {
	service, _, responsesGateway := createService(t)
	response, err := responsesGateway.Create("You are a bot", "Hi", "https://example.com", "Hello", "gpt-11-max", "text-embedding-medium", 0.5)
	require.NoError(t, err)

	result, found := service.Find(response.Id)

	assert.True(t, found)
	assert.Equal(t, response.Id, result.Id)
	assert.Equal(t, "You are a bot", result.SystemPrompt)
	assert.Equal(t, "Hi", result.UserQuery)
	assert.Equal(t, "https://example.com", result.Source)
	assert.Equal(t, "Hello", result.Response)
	assert.Equal(t, "gpt-11-max", result.ChatModel)
	assert.Equal(t, "text-embedding-medium", result.EmbeddingsModel)
	assert.Equal(t, float32(0.5), result.Temperature)
	assert.Nil(t, result.Score)
}

func TestScoredResponsesService_Find_NotFound(t *testing.T) {
	service, _, _ := createService(t)

	_, found := service.Find("bbaaaadd-7b55-4023-8c67-64204d30a900")

	assert.False(t, found)
}

func TestScoredResponsesService_List(t *testing.T) {
	service, gateway, responsesGateway := createService(t)
	response1, err := responsesGateway.Create("You are a bot", "Hi", "https://example.com", "Hello", "gpt-11-max", "text-embedding-medium", 0.5)
	require.NoError(t, err)
	response2, err := responsesGateway.Create("You are a bot", "Hi", "https://example.com", "Hello", "gpt-11-max", "text-embedding-medium", 0.5)
	require.NoError(t, err)
	_, err = gateway.Save(response1.Id, 11, 12, 13, 14)
	require.NoError(t, err)

	results, err := service.List()
	assert.NoError(t, err)

	assert.Len(t, results, 2)
	assert.Equal(t, []string{response2.Id, response1.Id}, slicesupport.Map(results, func(r scores.ResponseWithScore) string { return r.Id }))
	assert.Equal(t, "You are a bot", results[0].SystemPrompt)
	assert.Equal(t, "Hi", results[0].UserQuery)
	assert.Equal(t, "https://example.com", results[0].Source)
	assert.Equal(t, "Hello", results[0].Response)
	assert.Equal(t, "gpt-11-max", results[0].ChatModel)
	assert.Equal(t, "text-embedding-medium", results[0].EmbeddingsModel)
	assert.Equal(t, float32(0.5), results[0].Temperature)
	assert.Nil(t, results[0].Score)

	assert.Equal(t, 11, results[1].Score.Relevance)
	assert.Equal(t, 12, results[1].Score.Correctness)
	assert.Equal(t, 13, results[1].Score.AppropriateTone)
	assert.Equal(t, 14, results[1].Score.Politeness)
}

func createService(t *testing.T) (*scores.ScoredResponsesService, *scores.Gateway, *query.ResponsesGateway) {
	testDb := testsupport.NewTestDb(t)
	gateway := scores.NewGateway(testDb.DB)
	responsesGateway := query.NewResponsesGateway(testDb.DB)
	return scores.NewScoredResponsesService(gateway, responsesGateway), gateway, responsesGateway
}
