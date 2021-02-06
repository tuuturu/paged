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
		return errors.New("missing required DSN environment variable")
	}

	if c.Port == "" {
		return errors.New("missing required PORT environment variable")
	}

	return nil
}

func LoadConfig() (cfg *Config) {
	cfg = &Config{
		Port: "3000",
	}

	if port := os.Getenv("PAGED_LISTENING_PORT"); port != "" {
		cfg.Port = port
	}

	if dsn := os.Getenv("DSN"); dsn != "" {
		cfg.Database = parseDSN(dsn)
	}

	if discoveryURL, err := url.Parse(os.Getenv("DISCOVERY_URL")); err == nil {
		cfg.DiscoveryURL = discoveryURL
	}

	return cfg
}
