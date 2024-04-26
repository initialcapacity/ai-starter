package app

import (
	"context"
	"github.com/Azure/azure-sdk-for-go/sdk/ai/azopenai"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"log/slog"
	"net/http"
)

type model struct {
	Heading  string
	Label    string
	Query    string
	Response string
}

func Index() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = websupport.Render(w, Resources, "index", model{
			Heading: "What would you like to know?",
			Label:   "Query",
		})
	}
}

func Query(client *azopenai.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := r.ParseForm()
		if err != nil {
			slog.Error("unable to parse form", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		query := r.Form.Get("query")

		chatResponse, err := client.GetChatCompletions(context.Background(), azopenai.ChatCompletionsOptions{
			Messages: []azopenai.ChatRequestMessageClassification{
				&azopenai.ChatRequestSystemMessage{Content: to.Ptr("You are a reporter for a major world newspaper.")},
				&azopenai.ChatRequestSystemMessage{Content: to.Ptr("Write your response as if you were writing a short, high-quality news article for your paper. Limit your response to one paragraph.")},
				&azopenai.ChatRequestUserMessage{Content: azopenai.NewChatRequestUserMessageContent(query)},
			},
			DeploymentName: to.Ptr("gpt-4-turbo"),
		}, nil)
		if err != nil {
			slog.Error("unable fetch chat completion", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		response := chatResponse.ChatCompletions.Choices[0].Message.Content

		_ = websupport.Render(w, Resources, "index", model{
			Heading:  "What else would you like to know?",
			Label:    "New Query",
			Query:    query,
			Response: *response,
		})
	}
}
