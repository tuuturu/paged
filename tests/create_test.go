package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/tuuturu/pager-event-service/pkg/core/models"
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

			resultEvent := eventFromJSonBytes(result.Body.Bytes())
			assert.Assert(t, resultEvent.Id != "")

			resultEvent.Id = ""

			assert.Equal(t, tc.expectBody, resultEvent)
		})
	}
}

func TestEnsureTimestamp(t *testing.T) {
	testCases := []struct {
		name string

		withEvent models.Event
	}{
		{
			name: "Should have timestamp when not sending one",

			withEvent: models.Event{
				Title:       "Cool event",
				Description: "Semi interesting description",
			},
		},
		{
			name: "Should have timestamp when sending one",

			withEvent: models.Event{
				Title:       "Cool event",
				Description: "Semi interesting description",
				Timestamp:   "1612609653598734910",
			},
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

			id := createEvent(t, env, tc.withEvent)

			response, err := env.DoRequest(fmt.Sprintf("/events/%s", id), http.MethodGet, nil)
			assert.NilError(t, err)

			resultEvent := eventFromJSonBytes(response.Body.Bytes())

			assert.Assert(t, resultEvent.Timestamp != "")
		})
	}
}
