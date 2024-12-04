package evaluation

import "strconv"

type ScoreReporter struct {
}

func NewScoreReporter() ScoreReporter {
	return ScoreReporter{}
}

func (r ScoreReporter) Report(results chan ScoredResponse) chan []string {
	lines := make(chan []string)

	go func() {
		for result := range results {
			lines <- []string{
				result.Response.Query,
				result.Response.Response,
				result.Response.Source,
				strconv.Itoa(result.Score.Relevance),
				strconv.Itoa(result.Score.Correctness),
				strconv.Itoa(result.Score.AppropriateTone),
				strconv.Itoa(result.Score.Politeness),
			}
		}

		close(lines)
	}()

	return lines
}
