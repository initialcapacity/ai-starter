package websupport_test

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"testing"
)

func TestServer(t *testing.T) {
	server := websupport.NewServer(func(mux *http.ServeMux) {
		testsupport.Handle(mux, "GET /", "You passed the test")
	})

	port, _ := server.Start("localhost", 0)
	testsupport.AssertHealthy(t, port, "/")
	defer func(server *websupport.Server) {
		_ = server.Stop()
	}(server)

	response, err := http.Get(fmt.Sprintf("http://localhost:%d/", port))
	assert.NoError(t, err)

	body, err := io.ReadAll(response.Body)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "You passed the test")
}
