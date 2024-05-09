package feedsupport_test

import (
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestParser_AllLinks(t *testing.T) {
	endpoint, server := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.RssFeed(mux, "https://example.com")
	})
	defer testsupport.StopTestServer(t, server)
	parser := feedsupport.NewParser(http.Client{})

	links, err := parser.AllLinks(endpoint)

	assert.NoError(t, err)
	assert.Equal(t, []string{"https://example.com/pickles", "https://example.com/chicken"}, links)
}
