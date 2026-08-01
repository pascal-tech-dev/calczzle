package config

import "os"

// Config holds application runtime settings loaded from the environment.
type Config struct {
	Port string
}

// Load reads configuration from environment variables, applying defaults
// where values are unset.
func Load() Config {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return Config{Port: port}
}
