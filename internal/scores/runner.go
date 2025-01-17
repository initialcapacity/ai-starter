package scores

import (
	"github.com/initialcapacity/ai-starter/internal/query"
	"log/slog"
)

type Scorer interface {
	Score(response query.ChatResponse) (ResponseScore, error)
}

type ScoredResponse struct {
	Response query.ChatResponse
	Score    ResponseScore
}

type Runner struct {
	scorer Scorer
}

func NewRunner(scorer Scorer) Runner {
	return Runner{scorer: scorer}
}

func (r Runner) Score(responses chan query.ChatResponse) []ScoredResponse {
	scores := make([]ScoredResponse, 0)
	for response := range responses {
		score, err := r.scorer.Score(response)
		if err != nil {
			slog.Error("failed to score response",
				"query", response.Query,
				"response", response.Response,
				"err", err,
			)
		} else {
			scores = append(scores, ScoredResponse{response, score})
		}
	}
	return scores
}
