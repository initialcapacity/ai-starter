//go:build integration

package integration_test

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/initialcapacity/ai-starter/pkg/dbsupport"
	"github.com/initialcapacity/ai-starter/pkg/testsupport"
	"github.com/initialcapacity/ai-starter/pkg/websupport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func TestIntegration(t *testing.T) {
	client := http.Client{}
	openAiKey := websupport.RequireEnvironmentVariable[string]("OPEN_AI_KEY")
	dbUrl := "postgres://starter:starter@localhost:5432/starter_integration?sslmode=disable"

	createIntegrationDatabase(t)
	prepareIntegrationDatabase(t)
	prepareBuildDirectory(t)

	build(t, "migrate")
	build(t, "cannedrss")
	build(t, "collector")
	build(t, "analyzer")
	build(t, "app")

	testCtx, cancelCtx := context.WithCancel(context.Background())
	defer cancelCtx()

	startCommand(t, testCtx, "./build/cannedrss")
	startCommand(t, testCtx, "./build/app",
		fmt.Sprintf("DATABASE_URL=%s", dbUrl),
		fmt.Sprintf("OPEN_AI_KEY=%s", openAiKey),
		"HOST=localhost",
		"PORT=8234")

	runCommand(t, testCtx, "./build/migrate",
		fmt.Sprintf("DATABASE_URL=%s", dbUrl),
		"MIGRATIONS_LOCATION=file://../../databases/starter")
	runCommand(t, testCtx, "./build/collector",
		fmt.Sprintf("DATABASE_URL=%s", dbUrl),
		"FEEDS=http://localhost:8123")
	runCommand(t, testCtx, "./build/analyzer",
		fmt.Sprintf("DATABASE_URL=%s", dbUrl),
		fmt.Sprintf("OPEN_AI_KEY=%s", openAiKey))

	getResponse, err := client.Get("http://localhost:8234/")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, getResponse.StatusCode)
	getBody := readBody(t, getResponse)
	assert.Contains(t, getBody, "What would you like to know")

	postResponse, err := client.Post("http://localhost:8234/", "application/x-www-form-urlencoded", strings.NewReader("query=tell%20me%20about%20pickles"))
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, postResponse.StatusCode)
	postBody := readBody(t, postResponse)
	assert.Contains(t, postBody, "http://localhost:8123/pickles")
	assert.Contains(t, postBody, "</html>")
}

func createIntegrationDatabase(t *testing.T) {
	testsupport.WithSuperDb(t, func(superDb *sql.DB) {
		execute(t, superDb, "drop database if exists starter_integration")
		execute(t, superDb, "create database starter_integration")
		execute(t, superDb, "grant all privileges on database starter_integration to starter")
	})
}

func prepareIntegrationDatabase(t *testing.T) {
	integrationDb := dbsupport.CreateConnection("postgres://super_test@localhost:5432/starter_integration?sslmode=disable")
	defer func() {
		err := integrationDb.Close()
		require.NoError(t, err)
	}()

	execute(t, integrationDb, "create extension if not exists vector")
	execute(t, integrationDb, "grant usage, create on schema public to starter")
}

func execute(t *testing.T, db *sql.DB, command string) {
	_, err := db.Exec(command)
	require.NoError(t, err, fmt.Sprintf("unable to execute %s", command))
}

func prepareBuildDirectory(t *testing.T) {
	err := os.RemoveAll("../../build")
	require.NoError(t, err)
	err = os.MkdirAll("./build", os.ModePerm)
	require.NoError(t, err)
}

func build(t *testing.T, name string) {
	err := exec.Command("go", "build", "-o", fmt.Sprintf("./build/%s", name), fmt.Sprintf("../../cmd/%s", name)).Run()
	require.NoError(t, err, fmt.Sprintf("build %s failed", name))
}

func startCommand(t *testing.T, ctx context.Context, command string, environment ...string) {
	cmd := exec.CommandContext(ctx, command)
	cmd.Env = append(cmd.Env, environment...)
	err := cmd.Start()
	require.NoError(t, err, fmt.Sprintf("failed to start %s", command))
}

func runCommand(t *testing.T, ctx context.Context, command string, environment ...string) {
	cmd := exec.CommandContext(ctx, command)
	cmd.Env = append(cmd.Env, environment...)
	err := cmd.Run()
	require.NoError(t, err, fmt.Sprintf("failed to run %s", command))
}

func readBody(t *testing.T, response *http.Response) string {
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err, "unable to read response body")
	return string(body)
}
