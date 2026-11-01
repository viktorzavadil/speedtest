package config

import (
	"errors"
	"flag"
	"log"
	"os"
	"time"
)

type Config struct {
	TargetUrl             string
	ApiKey                string
	Streams               uint
	DownloadPhaseDuration time.Duration
	UploadPhaseDuration   time.Duration
	IdlePhaseDuration     time.Duration
	GracePeriod           time.Duration

	StatisticsInterval         time.Duration
	LatencyProbeInterval       time.Duration
	TargetMetricsProbeInterval time.Duration
}

func Load() (config *Config) {
	config = &Config{}

	apiKey := os.Getenv("API_KEY")
	if apiKey == "" {
		log.Fatal("API_KEY must be set")
	}
	config.ApiKey = apiKey

	flag.StringVar(&config.TargetUrl, "target", "", "target URL, e.g. http://target:8080")
	flag.UintVar(&config.Streams, "streams", 4, "number of parallel streams")
	flag.DurationVar(&config.IdlePhaseDuration, "idle-phase-duration", 5*time.Second, "duration of idle phase")
	flag.DurationVar(&config.DownloadPhaseDuration, "download-phase-duration", 10*time.Second, "duration of download phase")
	flag.DurationVar(&config.UploadPhaseDuration, "upload-phase-duration", 10*time.Second, "duration of upload phase")

	flag.Parse()

	config.GracePeriod = time.Duration(2) * time.Second

	config.StatisticsInterval = time.Duration(1) * time.Second
	config.LatencyProbeInterval = time.Duration(500) * time.Millisecond
	config.TargetMetricsProbeInterval = time.Duration(1000) * time.Millisecond

	err := validate(config)
	if err != nil {
		log.Fatal(err)
	}

	return config
}

func validate(config *Config) error {
	if config.TargetUrl == "" {
		return errors.New("--target is required")
	}
	if config.Streams < 1 {
		return errors.New("--streams must be greater than zero")
	}
	if config.IdlePhaseDuration < 0 {
		return errors.New("--idle-phase-duration must be greater than or equal to zero")
	}
	if config.DownloadPhaseDuration < 0 {
		return errors.New("--download-phase-duration must be greater than or equal to zero")
	}
	if config.UploadPhaseDuration < 0 {
		return errors.New("--upload-phase-duration must be greater than or equal to zero")
	}
	return nil
}
