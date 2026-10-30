package config

import (
	"flag"
	"log"
	"os"
	"time"
)

type Config struct {
	TargetUrl string
	ApiKey    string
	Duration  time.Duration
	Streams   uint

	StatisticsInterval   time.Duration
	LatencyProbeInterval time.Duration
}

func Load() (config *Config) {
	config = &Config{}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY must be set")
	}
	config.ApiKey = apiKey

	flag.StringVar(&config.TargetUrl, "target", "", "target URL, e.g. http://target:8080")
	flag.DurationVar(&config.Duration, "duration", 10*time.Second, "duration of each phase")
	flag.UintVar(&config.Streams, "streams", 4, "number of parallel streams")

	flag.Parse()

	config.StatisticsInterval = time.Duration(1) * time.Second
	config.LatencyProbeInterval = time.Duration(500) * time.Millisecond

	return config
}
