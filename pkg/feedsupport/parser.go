package feedsupport

import (
	"encoding/xml"
	"fmt"
	"net/http"
)

type Parser struct {
	client http.Client
}

func NewParser(client http.Client) Parser {
	return Parser{client: client}
}

func (p Parser) AllLinks(url string) ([]string, error) {
	response, err := p.client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("unable to fetch feed (%s): %w", url, err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	rss := Rss{}
	decoder := xml.NewDecoder(response.Body)
	err = decoder.Decode(&rss)
	if err != nil {
		return nil, fmt.Errorf("unable to decode feed (%s): %w", url, err)
	}

	var links []string
	for _, item := range rss.Channel.Items {
		links = append(links, item.Link)
	}

	return links, nil
}
