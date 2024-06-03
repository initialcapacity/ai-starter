package app

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"io/fs"
	"net/http"
)

func Handlers(aiClient ai.Client, db *sql.DB) func(mux *http.ServeMux) {
	embeddingsGateway := analyzer.NewEmbeddingsGateway(db)
	queryService := NewQueryService(embeddingsGateway, aiClient)

	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /", Index())
		mux.HandleFunc("POST /", Query(queryService))
		mux.HandleFunc("GET /health", Health)

		static, _ := fs.Sub(Resources, "resources/static")
		fileServer := http.FileServer(http.FS(static))
		mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	}
}
