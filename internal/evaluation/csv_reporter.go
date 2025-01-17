package evaluation

import (
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/initialcapacity/ai-starter/pkg/csvsupport"
	"log"
	"os"
	"strconv"
)

type CSVReporter struct {
}

func NewCSVReporter() CSVReporter {
	return CSVReporter{}
}

func (r CSVReporter) WriteToCSV(filename string, lines [][]string) error {
	rows := [][]string{{"Query", "Response", "Source", "Relevance", "Correctness", "Appropriate Tone", "Politeness"}}
	rows = append(rows, lines...)

	csvFile, err := os.Create(filename)
	defer func(csvFile *os.File) {
		_ = csvFile.Close()
	}(csvFile)
	if err != nil {
		log.Fatalln("failed to open file", err)
	}

	return csvsupport.WriteCSV(csvFile, rows)
}

func (r CSVReporter) Lines(results []scores.ScoredResponse) [][]string {
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
