package evaluation

import (
	"database/sql"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
)

type ScoresGateway struct {
	db *sql.DB
}

func NewScoresGateway(db *sql.DB) *ScoresGateway {
	return &ScoresGateway{db: db}
}

func (g *ScoresGateway) Save(queryResponseId string, relevance int, correctness int, appropriateTone int, politeness int) (string, error) {
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
