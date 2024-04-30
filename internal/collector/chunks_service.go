package collector

import (
	"errors"
	"github.com/initialcapacity/ai-starter/pkg/ai"
)

type ChunksService struct {
	chunker ai.Chunker
	gateway *ChunksGateway
}

func NewChunksService(chunker ai.Chunker, gateway *ChunksGateway) *ChunksService {
	return &ChunksService{chunker: chunker, gateway: gateway}
}

func (service ChunksService) SaveChunks(dataId, text string) error {
	chunks := service.chunker.Split(text)

	var saveErrors []error
	for _, chunk := range chunks {
		saveErrors = append(saveErrors, service.gateway.Save(dataId, chunk))
	}
	return errors.Join(saveErrors...)
}
