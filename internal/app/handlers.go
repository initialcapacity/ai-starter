package app

import (
	"github.com/initialcapacity/ai-starter/internal/ai"
	"io/fs"
	"net/http"
)

func Handlers(openAiKey string) func(mux *http.ServeMux) {
	aiClient := ai.NewClient(openAiKey)

	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /", Index())
		mux.HandleFunc("POST /", Query(aiClient))
		mux.HandleFunc("GET /health", Health)

		static, _ := fs.Sub(Resources, "resources/static")
		fileServer := http.FileServer(http.FS(static))
		mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	}
}
