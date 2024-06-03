package analyzer

import (
	"context"
	"errors"
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/collector"
	"log/slog"
)

type Analyzer struct {
	chunksGateway     *collector.ChunksGateway
	embeddingsGateway *EmbeddingsGateway
	embeddingCreator  embeddingCreator
}

func NewAnalyzer(chunksGateway *collector.ChunksGateway, embeddingsGateway *EmbeddingsGateway, embeddingCreator embeddingCreator) *Analyzer {
	return &Analyzer{chunksGateway: chunksGateway, embeddingsGateway: embeddingsGateway, embeddingCreator: embeddingCreator}
}

func (a *Analyzer) Analyze(ctx context.Context) error {
	slog.Info("Starting to analyze data")
	defer slog.Info("Finished analyzing data")

	ids, listErr := a.embeddingsGateway.UnprocessedIds()
	if listErr != nil {
		return fmt.Errorf("unable to list ids: %w", listErr)
	}

	slog.Info("found ids", "count", len(ids))
	var idErrors []error
	for _, id := range ids {
		record, err := a.chunksGateway.Get(id)
		if err != nil {
			idErrors = append(idErrors, fmt.Errorf("error getting content for id=%s: %w", id, err))
			continue
		}

		slog.Info("fetching embedding for", "id", id)
		embedding, err := a.embeddingCreator.CreateEmbedding(ctx, record.Content)
		if err != nil {
			idErrors = append(idErrors, fmt.Errorf("error fetching embedding for id=%s: %w", id, err))
			continue
		}

		slog.Info("saving embedding for", "id", id)
		err = a.embeddingsGateway.Save(id, embedding)
		if err != nil {
			idErrors = append(idErrors, fmt.Errorf("error saving embedding for id=%s: %w", id, err))
		}
	}

	return errors.Join(idErrors...)
}

type embeddingCreator interface {
	CreateEmbedding(ctx context.Context, text string) ([]float32, error)
}
