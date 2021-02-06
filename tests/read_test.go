package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/tuuturu/paged/pkg/core/models"
	"gotest.tools/assert"
)

func TestGetEvent(t *testing.T) {
	testCases := []struct {
		name string

		withEvent models.Event

		expectEqual  bool
		expectStatus int
	}{
		{
			name: "Should return correct data and code",

			withEvent: models.Event{
				Timestamp:   "123456789",
				Title:       "New event!",
				Description: "Has description",
				ImageURL:    "https://via.placeholder.com/150x150",
				ReadMoreURL: "https://news.tuuturu.org/hash",
			},

			expectEqual:  true,
			expectStatus: http.StatusOK,
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

			result, err := env.DoRequest(fmt.Sprintf("/events/%s", id), http.MethodGet, nil)
			assert.NilError(t, err)

			assert.Equal(t, tc.expectStatus, result.Code)

			resultEvent := eventFromJSonBytes(result.Body.Bytes())

			resultEvent.Id = ""

			assert.Assert(t, tc.withEvent == resultEvent, tc.expectEqual)
		})
	}
}

func TestGetMissingEvent(t *testing.T) {
	testCases := []struct {
		name string

		withID string

		expectCode int
	}{
		{
			name: "Should return 404",

			withID: "this-id-does-not-exist",

			expectCode: http.StatusNotFound,
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

			result, err := env.DoRequest(fmt.Sprintf("/events/%s", tc.withID), http.MethodGet, nil)
			assert.NilError(t, err)

			assert.Equal(t, tc.expectCode, result.Code)
		})
	}
}

func TestGetAllEvents(t *testing.T) {
	testCases := []struct {
		name string

		with []models.Event
	}{
		{
			name: "Should return a single event in a list when given one event",

			with: []models.Event{
				{
					Title:       "Awesome event",
					Description: "Very descriptive description",
				},
			},
		},
		{
			name: "Should return an empty list when given zero events",

			with: []models.Event{},
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

			createEvents(t, env, tc.with)

			response, err := env.DoRequest("/events", http.MethodGet, nil)
			assert.NilError(t, err)

			events := make([]models.Event, 0)

			err = json.Unmarshal(response.Body.Bytes(), &events)
			assert.NilError(t, err)
			assert.Assert(t, !bytes.Equal(response.Body.Bytes(), []byte("null")))

			assert.Equal(t, len(tc.with), len(events))
		})
	}
}

func TestEmptyEvents(t *testing.T) {
	env, err := CreateTestEnvironment()
	assert.NilError(t, err)

	defer func() {
		_ = env.Teardown()
	}()

	response, err := env.DoRequest("/events", http.MethodGet, nil)
	assert.NilError(t, err)

	assert.Assert(t, bytes.Equal([]byte("[]"), response.Body.Bytes()))
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		name string

		withFilter string
		withEvents []models.Event

		expectResults int
	}{
		{
			name: "Should exclude nothing without a filter",

			withFilter: "",
			withEvents: []models.Event{
				{
					Title:       "First event",
					Description: "some descript",
				},
				{
					Title:       "Second event",
					Description: "other descript",
				},
			},

			expectResults: 2,
		},
		{
			name: "Should exclude unread events",

			withFilter: "read=true",
			withEvents: []models.Event{
				{
					Title:       "This is read",
					Description: "Already read the description",
					Read:        true,
				},
				{
					Title:       "This is unread",
					Description: "Havent yet read the description",
					Read:        false,
				},
			},

			expectResults: 1,
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

			createEvents(t, env, tc.withEvents)

			filter := ""
			if tc.withFilter != "" {
				filter = fmt.Sprintf("?%s", tc.withFilter)
			}

			response, err := env.DoRequest(fmt.Sprintf("/events%s", filter), http.MethodGet, nil)
			assert.NilError(t, err)

			var resultEvents []models.Event

			err = json.Unmarshal(response.Body.Bytes(), &resultEvents)
			assert.NilError(t, err)

			assert.Equal(t, tc.expectResults, len(resultEvents))
		})
	}
}
