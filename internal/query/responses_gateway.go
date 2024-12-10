package query

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"time"
)

type ResponseRecord struct {
	Id           string
	SystemPrompt string
	UserQuery    string
	Source       string
	Response     string
	Model        string
	Temperature  float32
	CreatedAt    time.Time
}

type ResponsesGateway struct {
	db *sql.DB
}

func NewResponsesGateway(db *sql.DB) *ResponsesGateway {
	return &ResponsesGateway{db: db}
}

func (g *ResponsesGateway) Create(systemPrompt string, userQuery string, source string, response string, model string, temperature float32) (ResponseRecord, error) {
	return dbsupport.QueryOne(
		g.db,
		`insert into query_responses (system_prompt, user_query, source, response, model, temperature)
				values ($1, $2, $3, $4, $5, $6)
				returning id, system_prompt, user_query, source, response, model, temperature, created_at`,
		func(row *sql.Row, record *ResponseRecord) error {
			return row.Scan(&record.Id, &record.SystemPrompt, &record.UserQuery, &record.Source, &record.Response, &record.Model, &record.Temperature, &record.CreatedAt)
		},
		systemPrompt, userQuery, source, response, model, temperature,
	)
}

func (g *ResponsesGateway) List() ([]ResponseRecord, error) {
	return dbsupport.Query(
		g.db,
		`select id, system_prompt, user_query, source, response, model, temperature, created_at 
			from query_responses
			order by created_at desc`,
		func(row *sql.Rows, record *ResponseRecord) error {
			return row.Scan(&record.Id, &record.SystemPrompt, &record.UserQuery, &record.Source, &record.Response, &record.Model, &record.Temperature, &record.CreatedAt)
		})
}
