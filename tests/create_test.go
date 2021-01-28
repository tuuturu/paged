package tests

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/tuuturu/paged/pkg/core/models"
	"gotest.tools/assert"
)

func TestCreateEvent(t *testing.T) {
	testCases := []struct {
		name string

		with         models.Event
		expectBody   models.Event
		expectStatus int
	}{
		{
			name: "Should return correct status code and body upon success",

			with: models.Event{
				Timestamp:   "123456789",
				Title:       "New event!",
				Description: "Has description",
				ImageURL:    "https://via.placeholder.com/150x150",
				ReadMoreURL: "https://news.tuuturu.org/hash",
			},
			expectBody: models.Event{
				Timestamp:   "123456789",
				Title:       "New event!",
				Description: "Has description",
				ImageURL:    "https://via.placeholder.com/150x150",
				ReadMoreURL: "https://news.tuuturu.org/hash",
			},
			expectStatus: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			env, err := CreateTestEnvironment()
			assert.NilError(t, err)

			defer func() {
				_ = env.Teardown()
			}()

			result, err := env.DoRequest("/events", http.MethodPost, eventAsJSONBytes(tc.with))
			assert.NilError(t, err)

			assert.Equal(t, tc.expectStatus, result.Code)

			resultEvent := models.Event{}

			err = json.Unmarshal(result.Body.Bytes(), &resultEvent)
			assert.NilError(t, err)

			assert.Assert(t, resultEvent.Id != "")

			resultEvent.Id = ""

			assert.Equal(t, tc.expectBody, resultEvent)
		})
	}
}
