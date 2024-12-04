package evaluation

import (
	"context"
	"github.com/initialcapacity/ai-starter/internal/query"
	"log/slog"
	"strings"
	"sync"
)

type ChatResponse struct {
	Query    string
	Response string
	Source   string
}

type ChatResponseRetriever struct {
	queryService *query.Service
}

func NewChatResponseRetriever(queryService *query.Service) *ChatResponseRetriever {
	return &ChatResponseRetriever{queryService}
}

func (retriever ChatResponseRetriever) Retrieve(queries []string) chan ChatResponse {
	responses := make(chan ChatResponse)
	ctx := context.Background()
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(queries))

	go func() {
		waitGroup.Wait()
		close(responses)
	}()

	for _, q := range queries {
		go func() {
			defer waitGroup.Done()

			result, err := retriever.queryService.FetchResponse(ctx, q)
			if err != nil {
				slog.Error("failed to fetch response", "query", q, "err", err)
				return
			}

			response := strings.Builder{}
			for part := range result.Response {
				response.WriteString(part)
			}

			responses <- ChatResponse{
				Query:    q,
				Response: response.String(),
				Source:   result.Source,
			}
		}()
	}

	return responses
}
