package core

import "os"

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

	return cfg
}
