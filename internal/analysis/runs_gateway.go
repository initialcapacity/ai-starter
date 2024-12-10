package analysis

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"time"
)

type RunRecord struct {
	Id                string
	ChunksAnalyzed    int
	EmbeddingsCreated int
	NumberOfErrors    int
	CreatedAt         time.Time
}

type RunsGateway struct {
	db *sql.DB
}

func NewAnalysisRunsGateway(db *sql.DB) *RunsGateway {
	return &RunsGateway{db: db}
}

func (g RunsGateway) Create(chunksAnalyzed int, embeddingsCreated int, numberOfErrors int) (RunRecord, error) {
	return dbsupport.QueryOne(
		g.db,
		`insert into analysis_runs (chunks_analyzed, embeddings_created, errors)
				values ($1, $2, $3)
				returning id, chunks_analyzed, embeddings_created, errors, created_at`,
		func(row *sql.Row, record *RunRecord) error {
			return row.Scan(&record.Id, &record.ChunksAnalyzed, &record.EmbeddingsCreated, &record.NumberOfErrors, &record.CreatedAt)
		},
		chunksAnalyzed, embeddingsCreated, numberOfErrors,
	)
}

func (g RunsGateway) List() ([]RunRecord, error) {
	return dbsupport.Query(
		g.db,
		`select id, chunks_analyzed, embeddings_created, errors, created_at
			from analysis_runs
			order by created_at desc`,
		func(row *sql.Rows, record *RunRecord) error {
			return row.Scan(&record.Id, &record.ChunksAnalyzed, &record.EmbeddingsCreated, &record.NumberOfErrors, &record.CreatedAt)
		})
}
