package app

import (
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log/slog"
	"net/http"
)

func CollectionRuns(gateway *collection.RunsGateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runs, err := gateway.List()
		if err != nil {
			slog.Error("Could not list collection runs", "err", err)
			w.WriteHeader(500)
			return
		}

		_ = websupport.Render(w, Resources, "collection_runs", collectionRunsModel{runs})
	}
}

type collectionRunsModel struct {
	CollectionRuns []collection.RunRecord
}

func AnalysisRuns(gateway *analysis.RunsGateway) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		runs, err := gateway.List()
		if err != nil {
			slog.Error("Could not list analysis runs", "err", err)
			w.WriteHeader(500)
			return
		}

		_ = websupport.Render(w, Resources, "analysis_runs", analysisRunsModel{runs})
	}
}

type analysisRunsModel struct {
	AnalysisRuns []analysis.RunRecord
}
