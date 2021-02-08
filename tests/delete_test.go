package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/tuuturu/pager-event-service/pkg/core/models"
	"gotest.tools/assert"
)

func TestDeleteEvent(t *testing.T) {
	testCases := []struct {
		name string

		withEvent *models.Event

		expectStatusCode int
	}{
		{
			name: "Should work",

			withEvent: &models.Event{
				Timestamp:   "123456789",
				Title:       "New event!",
				Description: "Has description",
				ImageURL:    "https://via.placeholder.com/150x150",
				ReadMoreURL: "https://news.tuuturu.org/hash",
			},

			expectStatusCode: http.StatusNoContent,
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

			var id string
			if tc.withEvent != nil {
				id = createEvent(t, env, *tc.withEvent)
			}

			result, err := env.DoRequest(fmt.Sprintf("/events/%s", id), http.MethodDelete, nil)
			assert.NilError(t, err)

			assert.Equal(t, tc.expectStatusCode, result.Code)

			result, err = env.DoRequest(fmt.Sprintf("/events/%s", id), http.MethodDelete, nil)
			assert.NilError(t, err)

			assert.Equal(t, http.StatusNotFound, result.Code)
		})
	}
}

func TestDeleteMissingEvent(t *testing.T) {
	testCases := []struct {
		name string

		withID           string
		expectStatusCode int
	}{
		{
			name: "Should work",

			withID:           "non-existant-id",
			expectStatusCode: http.StatusNotFound,
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

			result, err := env.DoRequest(fmt.Sprintf("/events/%s", tc.withID), http.MethodDelete, nil)
			assert.NilError(t, err)

			assert.Equal(t, tc.expectStatusCode, result.Code)
		})
	}
}
