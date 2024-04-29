package dbsupport

import "database/sql"

func Query[T any](db *sql.DB, query string, mapping func(rows *sql.Rows, record *T) error, args ...any) ([]T, error) {
	var results []T
	rows, err := db.Query(query, args...)
	if err != nil {
		return results, err
	}
	defer func(rows *sql.Rows) {
		_ = rows.Close()
	}(rows)

	for rows.Next() {
		var record T
		err = mapping(rows, &record)
		if err != nil {
			return results, err
		}

		results = append(results, record)
	}

	return results, err
}

func QueryOne[T any](db *sql.DB, query string, mapping func(row *sql.Row, record *T) error, args ...any) (T, error) {
	var result T
	row := db.QueryRow(query, args...)
	err := mapping(row, &result)

	return result, err
}
