package main

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log"
	"net/http"
)

func main() {
	host := websupport.EnvironmentVariable("HOST", "localhost")
	port := websupport.EnvironmentVariable("PORT", 8123)

	server := websupport.NewServer(func(mux *http.ServeMux) {
		testsupport.RssFeed(mux, fmt.Sprintf("http://%s:%d", host, port))
		testsupport.Articles(mux)
	})

	_, done := server.Start(host, port)
	log.Fatal(<-done)
}
