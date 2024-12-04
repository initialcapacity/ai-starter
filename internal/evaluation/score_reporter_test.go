package evaluation_test

import (
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestScoreReporter_Report(t *testing.T) {
	reporter := evaluation.NewScoreReporter()
	scoringResponses := make(chan evaluation.ScoredResponse)

	lines := reporter.Report(scoringResponses)
	scoringResponses <- evaluation.ScoredResponse{
		Response: evaluation.ChatResponse{Query: "What's up?", Response: "Nothing much", Source: "https://example.com"},
		Score:    evaluation.ResponseScore{Relevance: 15, Correctness: 25, AppropriateTone: 35, Politeness: 45},
	}
	line := <-lines

	assert.Equal(t, []string{"What's up?", "Nothing much", "https://example.com", "15", "25", "35", "45"}, line)
}