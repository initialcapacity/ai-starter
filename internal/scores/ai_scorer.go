package scores

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/pkg/ai"
)

type aiClient interface {
	GetJsonChatCompletion(ctx context.Context, messages []ai.ChatMessage, schemaName string, schemaDescription string, jsonSchema string) (string, error)
}

type AiScorer struct {
	aiClient aiClient
}

func NewAiScorer(aiClient aiClient) AiScorer {
	return AiScorer{aiClient: aiClient}
}

func (s AiScorer) Score(response query.ChatResponse) (score ResponseScore, err error) {
	scoreResponse, err := s.aiClient.GetJsonChatCompletion(context.Background(), []ai.ChatMessage{
		{Role: ai.System, Content: fmt.Sprintf(`
			You are an expert QA professional. Below is a user's query about technology news, along with an assistant's response.
			Your task is to score the response on an integer scale from 0 to 100 on each of the following criteria:
			- Relevance: How relevant is the response to the user's query?'
			- Correctness: Does the response correctly answer or address the user's query?'
			- AppropriateTone: Does the response use appropriate tone (should be the tone of a tech journalist)?
			- Politeness: Does the response use polite language?\
		
			The JSON format of the response will be provided to you.
		
			Query: %s
		
			Response: %s
		`, response.Query, response.Response)},
	}, "ResponseScore", "The score of the response. Each number should be an integer between 0 and 100, inclusive",
		`{
			"type": "object",
			"properties": {
				"Relevance": {"type":  "integer"},
				"Correctness": {"type":  "integer"},
				"AppropriateTone": {"type":  "integer"},
				"Politeness": {"type":  "integer"}
			},
			"required": ["Relevance", "Correctness", "AppropriateTone", "Politeness"],
			"additionalProperties": false
		}`,
	)
	if err != nil {
		return score, err
	}

	err = json.Unmarshal([]byte(scoreResponse), &score)
	return score, err
}
