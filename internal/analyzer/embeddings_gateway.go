package analyzer

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/pgvector/pgvector-go"
)

type CitedChunkRecord struct {
	Content string
	Source  string
}

type EmbeddingsGateway struct {
	db *sql.DB
}

func NewEmbeddingsGateway(db *sql.DB) *EmbeddingsGateway {
	return &EmbeddingsGateway{db: db}
}

func (g *EmbeddingsGateway) UnprocessedIds() ([]string, error) {
	return dbsupport.Query(
		g.db,
		`select chunks.id from chunks
			left join public.embeddings e on chunks.id = e.chunk_id
			where e.id is null`,
		func(rows *sql.Rows, id *string) error { return rows.Scan(id) })
}

func (g *EmbeddingsGateway) Save(chunkId string, vector []float32) error {
	_, err := g.db.Exec("insert into embeddings (chunk_id, embedding) values ($1, $2)", chunkId, pgvector.NewVector(vector))
	return err
}

func (g *EmbeddingsGateway) FindSimilar(embedding []float32) (CitedChunkRecord, error) {
	return dbsupport.QueryOne(
		g.db,
		`select c.content, d.source
			from embeddings e
			join chunks c on c.id = e.chunk_id
			join data d on d.id = c.data_id
			order by e.embedding <=> $1 limit 1`,
		func(row *sql.Row, record *CitedChunkRecord) error {
			return row.Scan(&record.Content, &record.Source)
		},
		pgvector.NewVector(embedding),
	)
}
