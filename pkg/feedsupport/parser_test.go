package feedsupport_test

import (
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestParser_AllLinks(t *testing.T) {
	endpoint := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.HandleRssFeed(mux, "https://example.com")
	})
	parser := feedsupport.NewParser(http.Client{})

	links, err := parser.AllLinks(endpoint)

	assert.NoError(t, err)
	assert.Equal(t, []string{"https://example.com/pickles", "https://example.com/chicken"}, links)
}
