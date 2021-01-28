package tests

import (
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
			name: "Should work",

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

			id := createEvent(t, env, tc.withEvent)

			result, err := env.DoRequest(fmt.Sprintf("/events/%s", id), http.MethodGet, nil)
			assert.NilError(t, err)

			assert.Equal(t, tc.expectStatus, result.Code)

			resultEvent := models.Event{}

			err = json.Unmarshal(result.Body.Bytes(), &resultEvent)
			assert.NilError(t, err)

			resultEvent.Id = ""

			assert.Assert(t, tc.withEvent == resultEvent, tc.expectEqual)
		})
	}
}
