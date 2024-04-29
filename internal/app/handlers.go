package app

import (
	"github.com/initialcapacity/ai-starter/internal/ai"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"io/fs"
	"net/http"
)

func Handlers(openAiKey, databaseUrl string) func(mux *http.ServeMux) {
	aiClient := ai.NewClient(openAiKey)
	db := dbsupport.CreateConnection(databaseUrl)
	embeddingsGateway := analyzer.NewEmbeddingsGateway(db)

	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /", Index())
		mux.HandleFunc("POST /", Query(aiClient, embeddingsGateway))
		mux.HandleFunc("GET /health", Health)

		static, _ := fs.Sub(Resources, "resources/static")
		fileServer := http.FileServer(http.FS(static))
		mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	}
}
