package analyzer

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/pgvector/pgvector-go"
)

type EmbeddingsGateway struct {
	db *sql.DB
}

func NewEmbeddingsGateway(db *sql.DB) *EmbeddingsGateway {
	return &EmbeddingsGateway{db: db}
}

func (g *EmbeddingsGateway) Save(dataId string, vector []float32) error {
	_, err := g.db.Exec("insert into embeddings (data_id, embedding) values ($1, $2)", dataId, pgvector.NewVector(vector))
	return err
}

func (g *EmbeddingsGateway) FindSimilar(embedding []float32) (string, error) {
	return dbsupport.QueryOne(g.db, "select data_id from embeddings order by embedding <=> $1 limit 1", func(row *sql.Row, record *string) error {
		return row.Scan(record)
	}, pgvector.NewVector(embedding))
}
