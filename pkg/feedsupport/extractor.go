package feedsupport

import (
	"fmt"
	"golang.org/x/net/html"
	"net/http"
	"strings"
)

type Extractor struct {
	client http.Client
}

func NewExtractor(client http.Client) Extractor {
	return Extractor{client: client}
}

func (e Extractor) FullText(url string) (string, error) {
	response, err := e.client.Get(url)
	if err != nil {
		return "", fmt.Errorf("unable to fetch url (%s): %w", url, err)
	}
	defer func() {
		_ = response.Body.Close()
	}()

	var result strings.Builder
	document := html.NewTokenizer(response.Body)
	previousToken := document.Token()
tokenLoop:
	for {
		tokenType := document.Next()
		switch tokenType {
		case html.ErrorToken:
			break tokenLoop
		case html.StartTagToken:
			previousToken = document.Token()
		case html.TextToken:
			if previousToken.Data == "script" || previousToken.Data == "style" {
				continue
			}
			textContent := strings.TrimSpace(html.UnescapeString(string(document.Text())))
			if len(textContent) > 0 {
				result.WriteString(textContent)
			}
		default:
			continue
		}
	}

	return result.String(), nil
}
