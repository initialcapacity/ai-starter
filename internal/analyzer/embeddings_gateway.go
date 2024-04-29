package analyzer

import (
	"database/sql"
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
