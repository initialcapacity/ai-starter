package app

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"github.com/initialcapacity/ai-starter/pkg/deferrable"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log/slog"
	"net/http"
)

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = websupport.Render(w, Resources, "index", nil)
	}
}

type model struct {
	Query    string
	Response deferrable.Deferrable[string]
	Source   string
}

func Query(aiClient ai.Client, embeddingsGateway *analyzer.EmbeddingsGateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("unable to parse form", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		query := r.Form.Get("query")
		embedding, err := aiClient.CreateEmbedding(r.Context(), query)
		if err != nil {
			slog.Error("unable to create embedding", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		record, err := embeddingsGateway.FindSimilar(embedding)
		if err != nil {
			slog.Error("unable to find similar embedding", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		responseChannel, err := aiClient.GetChatCompletion(r.Context(), []ai.ChatMessage{
			{Role: ai.System, Content: "You are a reporter for a major world newspaper."},
			{Role: ai.System, Content: "Write your response as if you were writing a short, high-quality news article for your paper. Limit your response to one paragraph."},
			{Role: ai.System, Content: fmt.Sprintf("Use the following article for context: %s", record.Content)},
			{Role: ai.User, Content: query},
		})
		if err != nil {
			slog.Error("unable fetch chat completion", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = websupport.Render(w, Resources, "response", model{
			Query:    query,
			Response: deferrable.New(w, responseChannel),
			Source:   record.Source,
		})
	}
}
