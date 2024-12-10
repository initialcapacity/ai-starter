package collection

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"time"
)

type RunRecord struct {
	Id                string
	FeedsCollected    int
	ArticlesCollected int
	ChunksCollected   int
	NumberOfErrors    int
	CreatedAt         time.Time
}

type RunsGateway struct {
	db *sql.DB
}

func NewRunsGateway(db *sql.DB) *RunsGateway {
	return &RunsGateway{db: db}
}

func (g RunsGateway) Create(feedsCollected int, articlesCollected int, chunksCollected int, numberOfErrors int) (RunRecord, error) {
	return dbsupport.QueryOne(
		g.db,
		`insert into collection_runs (feeds_collected, articles_collected, chunks_collected, errors)
				values ($1, $2, $3, $4)
				returning id, feeds_collected, articles_collected, chunks_collected, errors, created_at`,
		func(row *sql.Row, record *RunRecord) error {
			return row.Scan(&record.Id, &record.FeedsCollected, &record.ArticlesCollected, &record.ChunksCollected, &record.NumberOfErrors, &record.CreatedAt)
		},
		feedsCollected, articlesCollected, chunksCollected, numberOfErrors,
	)
}

func (g RunsGateway) List() ([]RunRecord, error) {
	return dbsupport.Query(
		g.db,
		`select id, feeds_collected, articles_collected, chunks_collected, errors, created_at 
			from collection_runs
			order by created_at desc`,
		func(row *sql.Rows, record *RunRecord) error {
			return row.Scan(&record.Id, &record.FeedsCollected, &record.ArticlesCollected, &record.ChunksCollected, &record.NumberOfErrors, &record.CreatedAt)
		})
}
