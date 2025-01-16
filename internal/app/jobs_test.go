package app_test

import (
	"fmt"
	"github.com/initialcapacity/ai-starter/internal/analysis"
	"github.com/initialcapacity/ai-starter/internal/app"
	"github.com/initialcapacity/ai-starter/internal/collection"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
)

func TestCollectionRuns(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	appEndpoint := testsupport.StartTestServer(t, app.Handlers(testsupport.NewTestAiClient(""), testDb.DB))

	gateway := collection.NewRunsGateway(testDb.DB)
	_, err := gateway.Create(34, 56, 78, 9)
	require.NoError(t, err)

	resp, err := http.Get(fmt.Sprintf("%s/jobs/collections", appEndpoint))
	require.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	body := string(bytes)
	assert.Contains(t, body, "34")
	assert.Contains(t, body, "56")
	assert.Contains(t, body, "78")
	assert.Contains(t, body, "9")
}

func TestAnalysisRuns(t *testing.T) {
	testDb := testsupport.NewTestDb(t)
	appEndpoint := testsupport.StartTestServer(t, app.Handlers(testsupport.NewTestAiClient(""), testDb.DB))

	gateway := analysis.NewRunsGateway(testDb.DB)
	_, err := gateway.Create(34, 56, 9)
	require.NoError(t, err)

	resp, err := http.Get(fmt.Sprintf("%s/jobs/analyses", appEndpoint))
	require.NoError(t, err)

	bytes, err := io.ReadAll(resp.Body)
	require.NoError(t, err)
	body := string(bytes)
	assert.Contains(t, body, "34")
	assert.Contains(t, body, "56")
	assert.Contains(t, body, "9")
}
