package evaluation_test

import (
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoreReporter_Report(t *testing.T) {
	reporter := evaluation.NewCSVReporter()

	lines := reporter.Lines([]scores.ScoredResponse{{
		Response: query.ChatResponse{Query: "What's up?", Response: "Nothing much", Source: "https://example.com"},
		Score:    scores.ResponseScore{Relevance: 15, Correctness: 25, AppropriateTone: 35, Politeness: 45},
	}})

	assert.Equal(t, [][]string{{"What's up?", "Nothing much", "https://example.com", "15", "25", "35", "45"}}, lines)
}
