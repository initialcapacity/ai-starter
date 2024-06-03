package app

import (
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

func Query(queryService *QueryService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("unable to parse form", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		query := r.Form.Get("query")
		result, err := queryService.FetchResponse(r.Context(), query)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		_ = websupport.Render(w, Resources, "response", model{
			Query:    query,
			Response: deferrable.New(w, result.Response),
			Source:   result.Source,
		})
	}
}
