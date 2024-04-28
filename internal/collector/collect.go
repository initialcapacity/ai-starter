package collector

import (
	"errors"
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"log/slog"
)

type Collector struct {
	extractor feedsupport.Extractor
	rssParser feedsupport.Parser
	gateway   *DataGateway
}

func New(rssParser feedsupport.Parser, extractor feedsupport.Extractor, gateway *DataGateway) *Collector {
	return &Collector{rssParser: rssParser, extractor: extractor, gateway: gateway}
}

func (c *Collector) Collect(feedUrls []string) error {
	slog.Info("Starting to collect data")
	defer slog.Info("Finished collecting data")

	links := c.getLinks(feedUrls)
	var linkErrors []error

	for _, link := range links {
		slog.Info("Found", "link", link)

		exists, err := c.gateway.Exists(link)
		if err != nil || exists {
			linkErrors = append(linkErrors, err)
			continue
		}

		text, err := c.extractor.FullText(link)
		if err != nil {
			linkErrors = append(linkErrors, err)
		}

		err = c.gateway.Save(link, text)
		if err != nil {
			linkErrors = append(linkErrors, err)
		}
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
