package collector

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
)

type DataGateway struct {
	db *sql.DB
}

func NewDataGateway(db *sql.DB) *DataGateway {
	return &DataGateway{db: db}
}

func (g *DataGateway) UnprocessedIds() ([]string, error) {
	return dbsupport.Query(
		g.db,
		`select data.id from data
			left join public.embeddings e on data.id = e.data_id
			where e.id is null`,
		func(rows *sql.Rows, id *string) error { return rows.Scan(id) })
}

func (g *DataGateway) GetContent(id string) (string, error) {
	return dbsupport.QueryOne(
		g.db,
		"select content from data where id = $1",
		func(row *sql.Row, content *string) error { return row.Scan(content) },
		id,
	)
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

func (g *DataGateway) Save(source, content string) error {
	_, err := g.db.Exec("insert into data (source, content) values ($1, $2)", source, content)
	return err
}
