package app

import (
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"io/fs"
	"log/slog"
	"net/http"
)

func Handlers(openAiKey string) func(mux *http.ServeMux) {
	keyCredential := azcore.NewKeyCredential(openAiKey)
	client, err := azopenai.NewClientForOpenAI("https://api.openai.com/v1", keyCredential, nil)
	if err != nil {
		slog.Error("unable to create client", err)
		panic("unable to create client")
	}

	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /", Index())
		mux.HandleFunc("POST /", Query(client))
		mux.HandleFunc("GET /health", Health)

		static, _ := fs.Sub(Resources, "resources/static")
		fileServer := http.FileServer(http.FS(static))
		mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	}
}
