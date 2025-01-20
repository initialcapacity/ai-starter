package evaluation

import (
	"errors"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"log/slog"
)

type PastResponseEvaluator struct {
	responsesGateway *query.ResponsesGateway
	scoresGateway    *scores.Gateway
	scorer           scores.Scorer
}

func NewPastResponseEvaluator(responsesGateway *query.ResponsesGateway, scoresGateway *scores.Gateway, scorer scores.Scorer) PastResponseEvaluator {
	return PastResponseEvaluator{
		responsesGateway: responsesGateway,
		scoresGateway:    scoresGateway,
		scorer:           scorer,
	}
}

func (l PastResponseEvaluator) Run() error {
	responses, err := l.responsesGateway.ListMissingScores()
	if err != nil {
		slog.Error("failed to list responses", "err", err)
		return err
	}

	slog.Info("found responses", "count", len(responses))
	successfulEvaluations := 0
	errs := make([]error, 0)
	for _, response := range responses {
		score, scoreErr := l.scorer.Score(query.ChatResponse{
			Query:    response.UserQuery,
			Response: response.Response,
			Source:   response.Source,
		})
		if scoreErr != nil {
			slog.Error("failed to score response", "response", response.Id)
			errs = append(errs, scoreErr)
			continue
		}

		_, saveErr := l.scoresGateway.Save(response.Id, score.Relevance, score.Correctness, score.AppropriateTone, score.Politeness)
		if saveErr != nil {
			slog.Error("failed to save response score", "response", response.Id)
			errs = append(errs, saveErr)
		} else {
			successfulEvaluations++
		}
	}

	slog.Info("successful evaluations", "count", successfulEvaluations, "errors", len(errs))
	return errors.Join(errs...)
}
