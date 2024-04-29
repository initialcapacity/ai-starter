package analyzer

import (
	"context"
	"errors"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/initialcapacity/ai-starter/internal/collector"
	"log/slog"
)

type Analyzer struct {
	dataGateway       *collector.DataGateway
	embeddingsGateway *EmbeddingsGateway
	client            *azopenai.Client
}

func NewAnalyzer(dataGateway *collector.DataGateway, embeddingsGateway *EmbeddingsGateway, client *azopenai.Client) *Analyzer {
	return &Analyzer{dataGateway: dataGateway, embeddingsGateway: embeddingsGateway, client: client}
}

func (a *Analyzer) Analyze(ctx context.Context) error {
	slog.Info("Starting to analyze data")
	defer slog.Info("Finished analyzing data")

	ids, listErr := a.dataGateway.UnprocessedIds()
	if listErr != nil {
		return fmt.Errorf("unable to list ids: %w", listErr)
	}

	slog.Info("found ids", "count", len(ids))
	var idErrors []error
	for _, id := range ids {
		text, err := a.dataGateway.GetContent(id)
		if err != nil {
			idErrors = append(idErrors, err)
			continue
		}

		slog.Info("fetching embedding for", "id", id)
		embeddings, err := a.client.GetEmbeddings(ctx, azopenai.EmbeddingsOptions{
			Input:          []string{text},
			DeploymentName: to.Ptr("text-embedding-3-large"),
		}, nil)
		if err != nil {
			idErrors = append(idErrors, err)
			continue
		}

		slog.Info("saving embedding for", "id", id)
		vector := embeddings.Data[0].Embedding
		err = a.embeddingsGateway.Save(id, vector)
		if err != nil {
			slog.Error("error saving embedding for", "id", id, "error", err)
			idErrors = append(idErrors, err)
		}
	}

	return errors.Join(idErrors...)
}
