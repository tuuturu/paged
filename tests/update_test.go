package tests

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/tuuturu/paged/pkg/core/models"
	"gotest.tools/assert"
)

func TestUpdateEvent(t *testing.T) {
	testCases := []struct {
		name string

		withEvent  models.Event
		withUpdate models.Event

		expectResult models.Event
		expectCode   int
	}{
		{
			name: "Should correctly update an existing event",

			withEvent: models.Event{
				Timestamp:   "123456789",
				Title:       "Boring title!",
				Description: "Boring description",
				ImageURL:    "https://via.placeholder.com/150x150",
				ReadMoreURL: "https://news.tuuturu.org/hash",
			},
			withUpdate: models.Event{
				Title:       "Awesome title!",
				Description: "Funny description with lots of interesting facts",
			},

			expectResult: models.Event{
				Timestamp:   "123456789",
				Title:       "Awesome title!",
				Description: "Funny description with lots of interesting facts",
				ImageURL:    "https://via.placeholder.com/150x150",
				ReadMoreURL: "https://news.tuuturu.org/hash",
			},
			expectCode: http.StatusOK,
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
			url := fmt.Sprintf("/events/%s", id)

			patchResult, err := env.DoRequest(url, http.MethodPatch, eventAsJSONBytes(tc.withUpdate))
			assert.NilError(t, err)

			assert.Equal(t, tc.expectCode, patchResult.Code)

			getResult, err := env.DoRequest(url, http.MethodGet, nil)
			assert.NilError(t, err)

			tc.expectResult.Id = id
			assert.Equal(t, tc.expectResult, eventFromJSonBytes(getResult.Body.Bytes()))
		})
	}
}

func TestUpdateMissingEvent(t *testing.T) {
	testCases := []struct {
		name string

		withID     string
		withUpdate models.Event

		expectCode int
	}{
		{
			name: "Should return 404 on PATCH missing event",

			withID: "this-id-is-bonkers",
			withUpdate: models.Event{
				Title: "This try should fail!",
			},

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

			result, err := env.DoRequest(fmt.Sprintf("/events/%s", tc.withID), http.MethodPatch, eventAsJSONBytes(tc.withUpdate))
			assert.NilError(t, err)

			assert.Equal(t, tc.expectCode, result.Code)
		})
	}
}
