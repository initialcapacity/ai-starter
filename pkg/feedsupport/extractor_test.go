package feedsupport_test

import (
	"github.com/initialcapacity/ai-starter/pkg/feedsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestExtractor_FullText(t *testing.T) {
	endpoint, server := testsupport.StartTestServer(t, func(mux *http.ServeMux) {
		testsupport.Handle(mux, "/", `
				<html lang="en">
					<body>
						<script>const ignoreMe = true</script>
						<style>.ignore {content: 'me too'}</style>
						<p>some text</p>
					</body>
				</html>
			`)
	})
	defer testsupport.StopTestServer(t, server)
	extractor := feedsupport.NewExtractor(http.Client{})

	text, err := extractor.FullText(endpoint)

	assert.NoError(t, err)
	assert.Equal(t, "some text", text)
}
