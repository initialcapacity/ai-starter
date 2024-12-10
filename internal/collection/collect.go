package collection

import (
	"errors"
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"log/slog"
)

type Collector struct {
	rssParser     feedsupport.Parser
	extractor     feedsupport.Extractor
	gateway       *DataGateway
	chunksService *ChunksService
	runsGateway   *RunsGateway
}

func New(rssParser feedsupport.Parser, extractor feedsupport.Extractor,
	gateway *DataGateway, chunksService *ChunksService, runsGateway *RunsGateway) *Collector {
	return &Collector{rssParser, extractor, gateway, chunksService, runsGateway}
}

func (c *Collector) Collect(feedUrls []string) error {
	slog.Info("Starting to collect data")
	defer slog.Info("Finished collecting data")

	links := c.getLinks(feedUrls)
	var linkErrors []error
	articlesProcessed := 0
	chunksProcessed := 0

	for _, link := range links {
		slog.Info("Found", "link", link)

		exists, err := c.gateway.Exists(link)
		if err != nil {
			linkErrors = append(linkErrors, err)
		}
		if err != nil || exists {
			continue
		}
		articlesProcessed += 1

		text, err := c.extractor.FullText(link)
		if err != nil {
			linkErrors = append(linkErrors, err)
			continue
		}

		dataId, err := c.gateway.Save(link, text)
		if err != nil {
			linkErrors = append(linkErrors, err)
			continue
		}

		chunks, err := c.chunksService.SaveChunks(dataId, text)
		if err != nil {
			linkErrors = append(linkErrors, err)
		}
		chunksProcessed += chunks
	}

	_, err := c.runsGateway.Create(len(feedUrls), articlesProcessed, chunksProcessed, len(linkErrors))
	if err != nil {
		linkErrors = append(linkErrors, err)
	}

	return errors.Join(linkErrors...)
}

func (c *Collector) getLinks(feedUrls []string) []string {
	var allLinks []string
	for _, url := range feedUrls {
		slog.Info("Collecting", "url", url)
		links, _ := c.rssParser.AllLinks(url)
		allLinks = append(allLinks, links...)
	}
	return allLinks
}
