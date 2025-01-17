package evaluation_test

import (
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarkdownReporter_Report(t *testing.T) {
	reporter := evaluation.NewMarkdownReporter()

	markdown := reporter.Report([]scores.ScoredResponse{
		{
			Response: query.ChatResponse{Query: "What's up?", Response: "Nothing much", Source: "https://example.com"},
			Score:    scores.ResponseScore{Relevance: 15, Correctness: 25, AppropriateTone: 35, Politeness: 45},
		},
	})

	assert.Equal(t, `# Evaluation Results

---

## Query

What's up?

## Response

Source: https://example.com

Nothing much

## Scores

| Relevance | Correctness | Appropriate Tone | Politeness |
| --------- | ----------- | ---------------- | ---------- |
| 15        | 25          | 35               | 45         |

---

`, markdown)
}
