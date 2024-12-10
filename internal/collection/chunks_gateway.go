package collection

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
)

type ChunkRecord struct {
	DataId  string
	Content string
}

type ChunksGateway struct {
	db *sql.DB
}

func NewChunksGateway(db *sql.DB) *ChunksGateway {
	return &ChunksGateway{db: db}
}

func (g *ChunksGateway) Get(id string) (ChunkRecord, error) {
	return dbsupport.QueryOne(
		g.db,
		"select data_id, content from chunks where id = $1",
		func(row *sql.Row, record *ChunkRecord) error { return row.Scan(&record.DataId, &record.Content) },
		id,
	)
}

func (g *ChunksGateway) Save(dataId, content string) error {
	_, err := g.db.Exec("insert into chunks (data_id, content) values ($1, $2)", dataId, content)
	return err
}
