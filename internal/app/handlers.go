package app

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/internal/scores"
	"github.com/initialcapacity/ai-starter/pkg/ai"
	"io/fs"
	"net/http"
)

func Handlers(aiClient ai.Client, db *sql.DB) func(mux *http.ServeMux) {
	collectionRunsGateway := collection.NewRunsGateway(db)
	analysisRunsGateway := analysis.NewRunsGateway(db)
	embeddingsGateway := analysis.NewEmbeddingsGateway(db)
	responsesGateway := query.NewResponsesGateway(db)
	scoresGateway := scores.NewGateway(db)

	queryService := query.NewService(embeddingsGateway, aiClient, responsesGateway)
	scoredResponsesService := scores.NewScoredResponsesService(scoresGateway, responsesGateway)

	return func(mux *http.ServeMux) {
		mux.HandleFunc("GET /", Index())
		mux.HandleFunc("POST /", Query(queryService))
		mux.HandleFunc("GET /health", Health)
		mux.HandleFunc("GET /jobs/collections", CollectionRuns(collectionRunsGateway))
		mux.HandleFunc("GET /jobs/analyses", AnalysisRuns(analysisRunsGateway))
		mux.HandleFunc("GET /query_responses", QueryResponses(scoredResponsesService))
		mux.HandleFunc("GET /query_responses/{id}", ShowQueryResponse(scoredResponsesService))

		static, _ := fs.Sub(Resources, "resources/static")
		fileServer := http.FileServer(http.FS(static))
		mux.Handle("GET /static/", http.StripPrefix("/static/", fileServer))
	}
}
