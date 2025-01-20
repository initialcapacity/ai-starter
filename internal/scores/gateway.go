package scores

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"time"
)

type ScoreRecord struct {
	Id              string
	QueryResponseId string
	Relevance       int
	Correctness     int
	AppropriateTone int
	Politeness      int
	CreatedAt       time.Time
}

type Gateway struct {
	db *sql.DB
}

func NewGateway(db *sql.DB) *Gateway {
	return &Gateway{db: db}
}

func (g *Gateway) Save(queryResponseId string, relevance int, correctness int, appropriateTone int, politeness int) (string, error) {
	return dbsupport.QueryOne(
		g.db,
		`insert into response_scores (query_response_id, score, score_version)
			values ($1, json_build_object(
				'relevance', $2::int,
				'correctness', $3::int,
				'appropriate_tone', $4::int,
				'politeness', $5::int
			), 1) returning id`,
		func(row *sql.Row, id *string) error { return row.Scan(id) },
		queryResponseId, relevance, correctness, appropriateTone, politeness,
	)
}

func (g *Gateway) FindForResponseId(responseId string) (ScoreRecord, error) {
	return dbsupport.QueryOne(
		g.db,
		`select id,
       		query_response_id,
       		(score -> 'relevance')::int as relevance,
			(score -> 'correctness')::int as correctness,
			(score -> 'appropriate_tone')::int as appropriate_tone,
			(score -> 'politeness')::int as politeness,
       		created_at
		from response_scores
		where query_response_id = $1 and score_version = 1
		limit 1`,
		func(row *sql.Row, record *ScoreRecord) error {
			return row.Scan(&record.Id, &record.QueryResponseId, &record.Relevance, &record.Correctness, &record.AppropriateTone, &record.Politeness, &record.CreatedAt)
		}, responseId)
}

func (g *Gateway) ListForResponseIds(responseIds []string) ([]ScoreRecord, error) {
	return dbsupport.Query(
		g.db,
		`select id,
       		query_response_id,
       		(score -> 'relevance')::int as relevance,
			(score -> 'correctness')::int as correctness,
			(score -> 'appropriate_tone')::int as appropriate_tone,
			(score -> 'politeness')::int as politeness,
       		created_at
		from response_scores
		where query_response_id = any($1) and score_version = 1
		order by created_at`,
		func(row *sql.Rows, record *ScoreRecord) error {
			return row.Scan(&record.Id, &record.QueryResponseId, &record.Relevance, &record.Correctness, &record.AppropriateTone, &record.Politeness, &record.CreatedAt)
		}, responseIds)
}
