package core

import (
	"errors"
	"net/url"
	"os"

	log "github.com/sirupsen/logrus"
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
		Port: getEnv("PORT", "3000"),
	}

	logLevel := getEnv("LOG_LEVEL", "error")
	switch logLevel {
	case "error":
		cfg.LogLevel = log.ErrorLevel
	default:
		cfg.LogLevel = log.InfoLevel
	}

	if dsn := os.Getenv("DSN"); dsn != "" {
		cfg.Database = parseDSN(dsn)
	}

	if discoveryURL, err := url.Parse(os.Getenv("DISCOVERY_URL")); err == nil {
		cfg.DiscoveryURL = discoveryURL
	}

	if clientID := os.Getenv("CLIENT_ID"); clientID != "" {
		cfg.ClientID = clientID
		cfg.ClientSecret = os.Getenv("CLIENT_SECRET")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	result := os.Getenv(key)

	if result == "" {
		return fallback
	}

	return result
}
