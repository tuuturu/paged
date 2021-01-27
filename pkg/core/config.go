package core

import (
	"errors"
	"net/url"
	"os"
)

func (c Config) Validate() error {
	if c.DiscoveryURL == nil {
		return errors.New("missing required DISCOVERY_URL environment variable")
	}

	if c.Database == nil {
		return errors.New("missing required DATABASE_URI environment variable")
	}

	return nil
}

func LoadConfig() (cfg *Config) {
	cfg = &Config{
		Port: "3000",
	}

	if databaseURI := os.Getenv("DATABASE_URI"); databaseURI != "" {
		cfg.Database = &DatabaseOptions{
			URI:          databaseURI,
			DatabaseName: os.Getenv("DATABASE_NAME"),
			Username:     os.Getenv("DATABASE_USERNAME"),
			Password:     os.Getenv("DATABASE_PASSWORD"),
		}
	}

	if discoveryURL, err := url.Parse(os.Getenv("DISCOVERY_URL")); err == nil {
		cfg.DiscoveryURL = discoveryURL
	}

	return cfg
}
