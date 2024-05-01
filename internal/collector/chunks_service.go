package collector

import (
	"errors"
)

type ChunksService struct {
	chunker Chunker
	gateway *ChunksGateway
}

func NewChunksService(chunker Chunker, gateway *ChunksGateway) *ChunksService {
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

type Chunker interface {
	Split(text string) []string
}
