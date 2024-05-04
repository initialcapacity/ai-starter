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
		testsupport.Handle(mux, "GET /", fmt.Sprintf(`
			<rss>
				<channel>
					<item><link>http://%s:%d/pickles</link></item>
					<item><link>http://%s:%d/chicken</link></item>
				</channel>
			</rss>
		`, host, port, host, port))

		testsupport.Handle(mux, "GET /chicken", "This is a page about chickens. Chickens have feathers and lay eggs.")
		testsupport.Handle(mux, "GET /pickles", "This is a page about pickles. Pickles are a green and salty snack.")
	})

	_, done := server.Start(host, port)
	log.Fatal(<-done)
}
