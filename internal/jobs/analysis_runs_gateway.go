package jobs

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"time"
)

type AnalysisRunRecord struct {
	Id                string
	ChunksAnalyzed    int
	EmbeddingsCreated int
	NumberOfErrors    int
	CreatedAt         time.Time
}

type AnalysisRunsGateway struct {
	db *sql.DB
}

func NewAnalysisRunsGateway(db *sql.DB) *AnalysisRunsGateway {
	return &AnalysisRunsGateway{db: db}
}

func (g AnalysisRunsGateway) Create(chunksAnalyzed int, embeddingsCreated int, numberOfErrors int) (AnalysisRunRecord, error) {
	return dbsupport.QueryOne(
		g.db,
		`insert into analysis_runs (chunks_analyzed, embeddings_created, errors)
				values ($1, $2, $3)
				returning id, chunks_analyzed, embeddings_created, errors, created_at`,
		func(row *sql.Row, record *AnalysisRunRecord) error {
			return row.Scan(&record.Id, &record.ChunksAnalyzed, &record.EmbeddingsCreated, &record.NumberOfErrors, &record.CreatedAt)
		},
		chunksAnalyzed, embeddingsCreated, numberOfErrors,
	)
}

func (g AnalysisRunsGateway) List() ([]AnalysisRunRecord, error) {
	return dbsupport.Query(
		g.db,
		`select id, chunks_analyzed, embeddings_created, errors, created_at
			from analysis_runs
			order by created_at desc`,
		func(row *sql.Rows, record *AnalysisRunRecord) error {
			return row.Scan(&record.Id, &record.ChunksAnalyzed, &record.EmbeddingsCreated, &record.NumberOfErrors, &record.CreatedAt)
		})
}
