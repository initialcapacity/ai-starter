package scores

import (
	"github.com/initialcapacity/ai-starter/internal/query"
	"github.com/initialcapacity/ai-starter/pkg/slicesupport"
	"log/slog"
	"time"
)

type ResponseWithScore struct {
	Id              string
	SystemPrompt    string
	UserQuery       string
	Source          string
	Response        string
	ChatModel       string
	EmbeddingsModel string
	Temperature     float32
	CreatedAt       time.Time
	Score           *ResponseScore
}

type ScoredResponsesService struct {
	gateway          *Gateway
	responsesGateway *query.ResponsesGateway
}

func NewScoredResponsesService(gateway *Gateway, responsesGateway *query.ResponsesGateway) *ScoredResponsesService {
	return &ScoredResponsesService{gateway: gateway, responsesGateway: responsesGateway}
}

func (s *ScoredResponsesService) Find(id string) (ResponseWithScore, bool) {
	response, err := s.responsesGateway.Find(id)
	if err != nil {
		slog.Error("unable to find query response", "id", id, "err", err)
		return ResponseWithScore{}, false
	}

	scoreRecord, err := s.gateway.FindForResponseId(id)
	var score *ResponseScore
	if err == nil {
		score = &ResponseScore{
			Relevance:       scoreRecord.Relevance,
			Correctness:     scoreRecord.Correctness,
			AppropriateTone: scoreRecord.AppropriateTone,
			Politeness:      scoreRecord.Politeness,
		}
	} else {
		slog.Debug("unable to find score", "queryResponseId", id, "err", err)
	}

	return ResponseWithScore{
		Id:              response.Id,
		SystemPrompt:    response.SystemPrompt,
		UserQuery:       response.UserQuery,
		Source:          response.Source,
		Response:        response.Response,
		ChatModel:       response.ChatModel,
		EmbeddingsModel: response.EmbeddingsModel,
		Temperature:     response.Temperature,
		CreatedAt:       response.CreatedAt,
		Score:           score,
	}, true
}

func (s *ScoredResponsesService) List() ([]ResponseWithScore, error) {
	responses, err := s.responsesGateway.List()
	if err != nil {
		slog.Error("unable to list responses", "err", err)
		return nil, err
	}

	ids := slicesupport.Map(responses, func(r query.ResponseRecord) string { return r.Id })
	scoreRecords, err := s.gateway.ListForResponseIds(ids)
	if err != nil {
		slog.Error("unable to list scores", "err", err)
		return nil, err
	}

	return slicesupport.Map(responses, toResponseWithScore(scoreRecords)), nil
}

func toResponseWithScore(scoreRecords []ScoreRecord) func(query.ResponseRecord) ResponseWithScore {
	return func(response query.ResponseRecord) ResponseWithScore {
		scoreRecord, found := slicesupport.Find(scoreRecords, func(r ScoreRecord) bool {
			return r.QueryResponseId == response.Id
		})
		var score *ResponseScore
		if found {
			score = &ResponseScore{
				Relevance:       scoreRecord.Relevance,
				Correctness:     scoreRecord.Correctness,
				AppropriateTone: scoreRecord.AppropriateTone,
				Politeness:      scoreRecord.Politeness,
			}
		}

		return ResponseWithScore{
			Id:              response.Id,
			SystemPrompt:    response.SystemPrompt,
			UserQuery:       response.UserQuery,
			Source:          response.Source,
			Response:        response.Response,
			ChatModel:       response.ChatModel,
			EmbeddingsModel: response.EmbeddingsModel,
			Temperature:     response.Temperature,
			CreatedAt:       response.CreatedAt,
			Score:           score,
		}
	}
}
