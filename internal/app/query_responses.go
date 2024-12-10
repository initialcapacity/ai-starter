package app

import (
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/pkg/slicesupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log/slog"
	"net/http"
	"time"
)

func QueryResponses(gateway *query.ResponsesGateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		records, err := gateway.List()
		if err != nil {
			slog.Error("Could not list query responses", "err", err)
			w.WriteHeader(500)
			return
		}

		_ = websupport.Render(w, Resources, "query_responses", queryResponsesModel{slicesupport.Map(records, recordToQueryResponse)})
	}
}

type queryResponsesModel struct {
	QueryResponses []QueryResponse
}

type QueryResponse struct {
	Id           string
	SystemPrompt string
	UserQuery    string
	Source       string
	Response     string
	Model        string
	Temperature  float32
	CreatedAt    time.Time
}

func recordToQueryResponse(record query.ResponseRecord) QueryResponse {
	return QueryResponse{
		Id:           record.Id,
		SystemPrompt: truncate(record.SystemPrompt, 100),
		UserQuery:    truncate(record.UserQuery, 100),
		Source:       truncate(record.Source, 100),
		Response:     truncate(record.Response, 100),
		Model:        record.Model,
		Temperature:  record.Temperature,
		CreatedAt:    record.CreatedAt,
	}
}

func truncate(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	} else {
		return text[:maxLength-3] + "..."
	}
}
