package app_test

import (
	"github.com/initialcapacity/ai-starter/internal/app"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"testing"
)

func TestHealth(t *testing.T) {
	server := websupport.NewServer(app.Handlers(""))
	port, _ := server.Start("localhost", 0)
	testsupport.AssertHealthy(t, port, "/health")
}
