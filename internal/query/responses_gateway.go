package query

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"time"
)

type ResponseRecord struct {
	Id              string
	SystemPrompt    string
	UserQuery       string
	Source          string
	Response        string
	ChatModel       string
	EmbeddingsModel string
	Temperature     float32
	CreatedAt       time.Time
}

type ResponsesGateway struct {
	db *sql.DB
}

func NewResponsesGateway(db *sql.DB) *ResponsesGateway {
	return &ResponsesGateway{db: db}
}

func (g *ResponsesGateway) Create(systemPrompt, userQuery, source, response, chatModel, embeddingsModel string, temperature float32) (ResponseRecord, error) {
	return dbsupport.QueryOne(
		g.db,
		`insert into query_responses (system_prompt, user_query, source, response, chat_model, embeddings_model, temperature)
				values ($1, $2, $3, $4, $5, $6, $7)
				returning id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature, created_at`,
		func(row *sql.Row, record *ResponseRecord) error {
			return row.Scan(&record.Id, &record.SystemPrompt, &record.UserQuery, &record.Source, &record.Response, &record.ChatModel, &record.EmbeddingsModel, &record.Temperature, &record.CreatedAt)
		},
		systemPrompt, userQuery, source, response, chatModel, embeddingsModel, temperature,
	)
}

func (g *ResponsesGateway) Find(id string) (ResponseRecord, error) {
	return dbsupport.QueryOne(
		g.db,
		`select id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature, created_at 
			from query_responses
			where id = $1`,
		func(row *sql.Row, record *ResponseRecord) error {
			return row.Scan(&record.Id, &record.SystemPrompt, &record.UserQuery, &record.Source, &record.Response, &record.ChatModel, &record.EmbeddingsModel, &record.Temperature, &record.CreatedAt)
		}, id)
}

func (g *ResponsesGateway) List() ([]ResponseRecord, error) {
	return dbsupport.Query(
		g.db,
		`select id, system_prompt, user_query, source, response, chat_model, embeddings_model, temperature, created_at 
			from query_responses
			order by created_at desc`,
		func(row *sql.Rows, record *ResponseRecord) error {
			return row.Scan(&record.Id, &record.SystemPrompt, &record.UserQuery, &record.Source, &record.Response, &record.ChatModel, &record.EmbeddingsModel, &record.Temperature, &record.CreatedAt)
		})
}

func (g *ResponsesGateway) ListMissingScores() ([]ResponseRecord, error) {
	return dbsupport.Query(
		g.db,
		`select r.id, r.system_prompt, r.user_query, r.source, r.response, r.chat_model, r.embeddings_model, r.temperature, r.created_at 
			from query_responses r
				left join response_scores e on e.query_response_id = r.id
				where e.id is null`,
		func(row *sql.Rows, record *ResponseRecord) error {
			return row.Scan(&record.Id, &record.SystemPrompt, &record.UserQuery, &record.Source, &record.Response, &record.ChatModel, &record.EmbeddingsModel, &record.Temperature, &record.CreatedAt)
		})
}
