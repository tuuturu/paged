package tests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/tuuturu/paged/pkg/core"
	"github.com/tuuturu/paged/pkg/core/models"
	"github.com/tuuturu/paged/pkg/core/router"
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

	cfg := &core.Config{
		DiscoveryURL: discoveryURL,
		Database: &core.DatabaseOptions{
			URI:          env.GetDatabaseBackendURI(),
			Username:     "postgres",
			Password:     dbPassword,
			DatabaseName: "postgres",
		},
	}

	err = cfg.Validate()
	if err != nil {
		return nil, fmt.Errorf("error validating config: %w", err)
	}

	env.TestServer = router.New(cfg)

	return env, nil
}

func eventAsJSONBytes(event models.Event) []byte {
	result, _ := json.Marshal(event)

	return result
}
