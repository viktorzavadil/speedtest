package main

import (
	"speedtest/target/internal/api"
	"speedtest/target/internal/config"
)

func main() {
	cfg := config.Load()
	server := api.NewServer(cfg)
	server.Serve()
}
