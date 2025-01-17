package evaluation

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"os"
	"strings"
)

type MarkdownReporter struct {
}

func NewMarkdownReporter() MarkdownReporter {
	return MarkdownReporter{}
}

func (r MarkdownReporter) WriteToFile(filename string, content string) error {
	file, err := os.Create(filename)
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		return err
	}

	_, err = file.WriteString(content)
	return err
}

func (r MarkdownReporter) Report(results []scores.ScoredResponse) string {
	builder := strings.Builder{}

	builder.WriteString(`# Evaluation Results

---
`)

	for _, result := range results {
		builder.WriteString(fmt.Sprintf(`
## Query

%s

## Response

Source: %s

%s

## Scores

| Relevance | Correctness | Appropriate Tone | Politeness |
| --------- | ----------- | ---------------- | ---------- |
| %d        | %d          | %d               | %d         |

---

`, result.Response.Query, result.Response.Source, result.Response.Response,
			result.Score.Relevance, result.Score.Correctness, result.Score.AppropriateTone, result.Score.Politeness))

	}

	return builder.String()
}
