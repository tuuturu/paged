package core

import (
	"net/url"
	"testing"

	"gotest.tools/assert"
)

func TestValidateConfig(t *testing.T) {
	testCases := []struct {
		name string

		withConfig Config

		expectFail bool
	}{
		{
			name: "Should not error when all required information is available",

			withConfig: Config{
				DiscoveryURL: &url.URL{
					Scheme: "https",
					Host:   "example.com",
				},
				Port: "5555",
				Database: &DatabaseOptions{
					URI:          "notadb:8833",
					DatabaseName: "postgres",
					Username:     "postgres",
					Password:     "postgres",
				},
			},

			expectFail: false,
		},
		{
			name: "Should error with no values",

			withConfig: Config{
				DiscoveryURL: nil,
				Port:         "",
				Database:     nil,
			},

			expectFail: true,
		},
		{
			name: "Should error with missing database and port",

			withConfig: Config{
				DiscoveryURL: &url.URL{
					Scheme: "https",
					Host:   "example.com",
				},
				Port:     "",
				Database: nil,
			},

			expectFail: true,
		},
		{
			name: "Should error with missing port",

			withConfig: Config{
				DiscoveryURL: &url.URL{
					Scheme: "https",
					Host:   "example.com",
				},
				Port: "",
				Database: &DatabaseOptions{
					URI:          "testuri:5432",
					DatabaseName: "postgres",
					Username:     "postgres",
					Password:     "postgres",
				},
			},

			expectFail: true,
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			err := tc.withConfig.Validate()

			assert.Equal(t, err != nil, tc.expectFail)
		})
	}
}
