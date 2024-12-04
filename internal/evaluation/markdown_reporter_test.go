package evaluation_test

import (
	"github.com/initialcapacity/ai-starter/internal/evaluation"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarkdownReporter_Report(t *testing.T) {
	reporter := evaluation.NewMarkdownReporter()

	markdown := reporter.Report([]evaluation.ScoredResponse{
		{
			Response: evaluation.ChatResponse{Query: "What's up?", Response: "Nothing much", Source: "https://example.com"},
			Score:    evaluation.ResponseScore{Relevance: 15, Correctness: 25, AppropriateTone: 35, Politeness: 45},
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
