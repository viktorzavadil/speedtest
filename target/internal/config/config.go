package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Port   uint16
	ApiKey string
}

func Load() (config *Config) {
	config = &Config{}

	portRaw := os.Getenv("PORT")
	if portRaw == "" {
		config.Port = 8080
	} else {
		port, err := strconv.ParseUint(portRaw, 10, 16)
		if err != nil {
			log.Fatal("failed to read PORT, number is expecting")
		}
		config.Port = uint16(port)
	}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY must be set")
	}
	config.ApiKey = apiKey

	return config
}
