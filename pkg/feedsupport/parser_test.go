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
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = w.Write([]byte(`
				<rss>
					<channel>
						<item><link>https://example.com/1</link></item>
						<item><link>https://example.com/2</link></item>
					</channel>
				</rss>
			`))
		})
	})
	defer testsupport.StopTestServer(t, server)
	parser := feedsupport.NewParser(http.Client{})

	links, err := parser.AllLinks(endpoint)

	assert.NoError(t, err)
	assert.Equal(t, []string{"https://example.com/1", "https://example.com/2"}, links)
}
