package evaluation

import (
	"errors"
	"fmt"
	"log/slog"
	"path"
)

type CannedResponseEvaluator struct {
	retriever   *ChatResponseRetriever
	scoreRunner ScoreRunner
	csvReporter CSVReporter
	mdReporter  MarkdownReporter
}

func NewCannedResponseEvaluator(retriever *ChatResponseRetriever, scoreRunner ScoreRunner, csvReporter CSVReporter, mdReporter MarkdownReporter) *CannedResponseEvaluator {
	return &CannedResponseEvaluator{retriever: retriever, scoreRunner: scoreRunner, csvReporter: csvReporter, mdReporter: mdReporter}
}

func (e CannedResponseEvaluator) Run(directory string, queries []string) error {
	results := e.retriever.Retrieve(queries)
	scores := e.scoreRunner.Score(results)
	if len(scores) == 0 {
		return errors.New("no scores were generated, there was likely a problem")
	}

	csvPath := path.Join(directory, "scores.csv")
	err := e.csvReporter.WriteToCSV(csvPath, e.csvReporter.Lines(scores))
	if err != nil {
		return fmt.Errorf("failed to write scores.csv: %w", err)
	}
	slog.Info("successfully wrote CSV", "file", csvPath)

	mdPath := path.Join(directory, "scores.md")
	err = e.mdReporter.WriteToFile(mdPath, e.mdReporter.Report(scores))
	if err != nil {
		return fmt.Errorf("failed to write scores.md: %w", err)
	}
	slog.Info("successfully wrote markdown", "file", mdPath)

	return nil
}
