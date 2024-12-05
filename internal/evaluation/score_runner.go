package evaluation

import "log/slog"

type Scorer interface {
	Score(response ChatResponse) (ResponseScore, error)
}

type ScoredResponse struct {
	Response ChatResponse
	Score    ResponseScore
}

type ScoreRunner struct {
	scorer Scorer
}

func NewScoreRunner(scorer Scorer) ScoreRunner {
	return ScoreRunner{scorer: scorer}
}

func (s ScoreRunner) Score(responses chan ChatResponse) []ScoredResponse {
	scores := make([]ScoredResponse, 0)
	for response := range responses {
		score, err := s.scorer.Score(response)
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
