package functions

import (
	"context"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/cloudevents/sdk-go/v2/event"
	"github.com/initialcapacity/ai-starter/internal/analyzer"
	"github.com/initialcapacity/ai-starter/internal/collector"
)

func init() {
	functions.CloudEvent("analyzer", triggerAnalyze)
	functions.CloudEvent("collector", triggerCollect)
}

func triggerCollect(ctx context.Context, e event.Event) error {
	return collector.Collect()
}

func triggerAnalyze(ctx context.Context, e event.Event) error {
	return analyzer.Analyze()
}
