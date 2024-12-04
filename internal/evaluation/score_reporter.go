package evaluation

import (
	"github.com/initialcapacity/ai-starter/pkg/csvsupport"
	"log"
	"os"
	"strconv"
)

type ScoreReporter struct {
}

func NewScoreReporter() ScoreReporter {
	return ScoreReporter{}
}

func (r ScoreReporter) WriteToCSV(filename string, lines [][]string) error {
	rows := [][]string{{"Query", "Response", "Source", "Relevance", "Correctness", "Appropriate Tone", "Politeness"}}
	rows = append(rows, lines...)

	csvFile, err := os.Create(filename)
	defer csvFile.Close()
	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	return csvsupport.WriteCSV(csvFile, rows)
}

func (r ScoreReporter) Report(results []ScoredResponse) [][]string {
	lines := make([][]string, 0)

	for _, result := range results {
		lines = append(lines, []string{
			result.Response.Query,
			result.Response.Response,
			result.Response.Source,
			strconv.Itoa(result.Score.Relevance),
			strconv.Itoa(result.Score.Correctness),
			strconv.Itoa(result.Score.AppropriateTone),
			strconv.Itoa(result.Score.Politeness),
		})
	}

	return lines
}
