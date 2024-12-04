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

func (s ScoreRunner) Score(responses chan ChatResponse) chan ScoredResponse {
	scores := make(chan ScoredResponse)
	go func() {
		for response := range responses {
			score, err := s.scorer.Score(response)
			if err != nil {
				slog.Error("failed to score response",
					"query", response.Query,
					"response", response.Response,
					"err", err,
				)
			} else {
				scores <- ScoredResponse{response, score}
			}
		}
		close(scores)
	}()
	return scores
}
