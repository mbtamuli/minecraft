package config

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Token         string `json:"DISCORD_TOKEN"`
	AllowedRoleID string `json:"ALLOWED_ROLE_ID"`
	ComposePath   string `json:"COMPOSE_FILE_PATH"`
}

func Load() *Config {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("FATAL: Could not read config.json: %v", err)
	}

	var cfg Config
	err = json.Unmarshal(file, &cfg)
	if err != nil {
		log.Fatalf("FATAL: Could not parse config.json: %v", err)
	}

	if cfg.Token == "" || cfg.AllowedRoleID == "" || cfg.ComposePath == "" {
		log.Fatal("FATAL: One or more required fields are missing in config.json.")
	}

	log.Printf("Loaded config: ComposePath=%s", cfg.ComposePath)
	return &cfg
}
