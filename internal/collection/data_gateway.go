package collection

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
)

type DataRecord struct {
	Content string
	Source  string
}

type DataGateway struct {
	db *sql.DB
}

func NewDataGateway(db *sql.DB) *DataGateway {
	return &DataGateway{db: db}
}

func (g *DataGateway) Exists(source string) (bool, error) {
	count, err := dbsupport.QueryOne(
		g.db,
		"select count(1) as count from data where source = $1",
		func(row *sql.Row, count *int) error { return row.Scan(count) },
		source,
	)

	return count > 0, err
}

func (g *DataGateway) Save(source, content string) (string, error) {
	return dbsupport.QueryOne(
		g.db,
		"insert into data (source, content) values ($1, $2) returning id",
		func(row *sql.Row, id *string) error { return row.Scan(id) },
		source, content,
	)
}
