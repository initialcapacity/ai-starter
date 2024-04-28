package collector

import "database/sql"

type DataGateway struct {
	db *sql.DB
}

func NewDataGateway(db *sql.DB) *DataGateway {
	return &DataGateway{db: db}
}

func (g *DataGateway) Exists(source string) (bool, error) {
	var count int
	row := g.db.QueryRow("select count(1) as count from data where source = $1", source)
	err := row.Scan(&count)
	return count > 0, err
}

func (g *DataGateway) Save(source, content string) error {
	_, err := g.db.Exec("insert into data (source, content) values ($1, $2)", source, content)
	return err
}
