package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/tuuturu/pager-event-service/pkg/core"
	"github.com/tuuturu/pager-event-service/pkg/core/models"
	"github.com/tuuturu/pager-event-service/pkg/core/router"
	"gotest.tools/assert"

	"github.com/oslokommune/go-gin-tools/pkg/v1/servicetesting"
	authtesting "github.com/oslokommune/go-oidc-middleware/pkg/v1/testing"
)

func CreateTestEnvironment() (*servicetesting.Environment, error) {
	authTestOptions := authtesting.NewTestTokenOptions()

	discoveryServer := authtesting.CreateTestDiscoveryServer(authTestOptions)
	bearerToken := authtesting.CreateTestToken(authTestOptions)

	discoveryURL, _ := url.Parse(discoveryServer.URL)

	dbPassword := "postgres"

	env, err := servicetesting.NewGinTestEnvironment(servicetesting.CreatePostgresDatabaseBackendOptions(dbPassword), bearerToken)
	if err != nil {
		return nil, fmt.Errorf("error creating test environment: %w", err)
	}

	parts := strings.Split(env.GetDatabaseBackendURI(), ":")
	dbURI := parts[0]
	dbPort := parts[1]

	cfg := &core.Config{
		DiscoveryURL: discoveryURL,
		Port:         "3000",
		Database: &core.DSN{
			Scheme:       "postgres",
			URI:          dbURI,
			Port:         dbPort,
			DatabaseName: "postgres",
			Username:     "postgres",
			Password:     dbPassword,
		},
	}

	err = cfg.Validate()
	if err != nil {
		_ = env.Teardown()

		return nil, fmt.Errorf("error validating config: %w", err)
	}

	env.TestServer = router.New(cfg)

	return env, nil
}

func eventAsJSONBytes(event models.Event) []byte {
	result, _ := json.Marshal(event)

	return result
}

func eventFromJSonBytes(buf []byte) (result models.Event) {
	_ = json.Unmarshal(buf, &result)

	return result
}

func createEvent(t *testing.T, env *servicetesting.Environment, event models.Event) string {
	result, err := env.DoRequest("/events", http.MethodPost, eventAsJSONBytes(event))
	assert.NilError(t, err)

	createdEvent := models.Event{}

	err = json.Unmarshal(result.Body.Bytes(), &createdEvent)
	assert.NilError(t, err)

	return createdEvent.Id
}

func createEvents(t *testing.T, env *servicetesting.Environment, events []models.Event) []string {
	ids := make([]string, len(events))

	for _, event := range events {
		id := createEvent(t, env, event)

		ids = append(ids, id)
	}

	return ids
}
