package testsupport

import (
	"fmt"
	"net/http"
)

func RssFeed(mux *http.ServeMux, articlesEndpoint string) {
	Handle(mux, "GET /", fmt.Sprintf(`
			<rss>
				<channel>
					<item><link>%s/pickles</link></item>
					<item><link>%s/chicken</link></item>
				</channel>
			</rss>
		`, articlesEndpoint, articlesEndpoint))
}

func Articles(mux *http.ServeMux) {
	Handle(mux, "GET /chicken", "This is a page about chickens. Chickens have feathers and lay eggs.")
	Handle(mux, "GET /pickles", "This is a page about pickles. Pickles are a green and salty snack.")
}
