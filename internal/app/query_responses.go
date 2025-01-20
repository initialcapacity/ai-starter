package app

import (
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/initialcapacity/ai-starter/pkg/slicesupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log/slog"
	"net/http"
	"time"
)

func QueryResponses(service *scores.ScoredResponsesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		responses, err := service.List()
		if err != nil {
			slog.Error("Could not list query responses", "err", err)
			w.WriteHeader(500)
			return
		}

		_ = websupport.Render(w, Resources, "query_responses", queryResponsesModel{slicesupport.Map(responses, truncateLongStrings)})
	}
}

func ShowQueryResponse(service *scores.ScoredResponsesService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.PathValue("id")
		response, found := service.Find(id)
		if !found {
			slog.Error("Could not find query response")
			w.WriteHeader(404)
			return
		}

		_ = websupport.Render(w, Resources, "query_response", showQueryResponseModel{response})
	}
}

type queryResponsesModel struct {
	Responses []scores.ResponseWithScore
}

type QueryResponse struct {
	Id              string
	SystemPrompt    string
	UserQuery       string
	Source          string
	Response        string
	ChatModel       string
	EmbeddingsModel string
	Temperature     float32
	CreatedAt       time.Time
}

func truncateLongStrings(response scores.ResponseWithScore) scores.ResponseWithScore {
	return scores.ResponseWithScore{
		Id:              response.Id,
		SystemPrompt:    truncate(response.SystemPrompt, 100),
		UserQuery:       truncate(response.UserQuery, 100),
		Source:          truncate(response.Source, 100),
		Response:        truncate(response.Response, 100),
		ChatModel:       response.ChatModel,
		EmbeddingsModel: response.EmbeddingsModel,
		Temperature:     response.Temperature,
		CreatedAt:       response.CreatedAt,
		Score:           response.Score,
	}
}

func truncate(text string, maxLength int) string {
	if len(text) <= maxLength {
		return text
	} else {
		return text[:maxLength-3] + "..."
	}
}

type showQueryResponseModel struct {
	Response scores.ResponseWithScore
}
