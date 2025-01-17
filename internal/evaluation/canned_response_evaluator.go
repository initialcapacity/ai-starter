package evaluation

import (
	"errors"
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"log/slog"
	"path"
)

type CannedResponseEvaluator struct {
	retriever   *query.ChatResponseRetriever
	scoreRunner scores.Runner
	csvReporter CSVReporter
	mdReporter  MarkdownReporter
}

func NewCannedResponseEvaluator(retriever *query.ChatResponseRetriever, scoreRunner scores.Runner, csvReporter CSVReporter, mdReporter MarkdownReporter) *CannedResponseEvaluator {
	return &CannedResponseEvaluator{retriever: retriever, scoreRunner: scoreRunner, csvReporter: csvReporter, mdReporter: mdReporter}
}

func (e CannedResponseEvaluator) Run(directory string, queries []string) error {
	results := e.retriever.Retrieve(queries)
	scoreList := e.scoreRunner.Score(results)
	if len(scoreList) == 0 {
		return errors.New("no scores were generated, there was likely a problem")
	}

	csvPath := path.Join(directory, "scores.csv")
	err := e.csvReporter.WriteToCSV(csvPath, e.csvReporter.Lines(scoreList))
	if err != nil {
		return fmt.Errorf("failed to write scores.csv: %w", err)
	}
	slog.Info("successfully wrote CSV", "file", csvPath)

	mdPath := path.Join(directory, "scores.md")
	err = e.mdReporter.WriteToFile(mdPath, e.mdReporter.Report(scoreList))
	if err != nil {
		return fmt.Errorf("failed to write scores.md: %w", err)
	}
	slog.Info("successfully wrote markdown", "file", mdPath)

	return nil
}
