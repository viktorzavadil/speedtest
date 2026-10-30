package main

import (
	"context"
	"os/signal"
	"speedtest/target/internal/api"
	"speedtest/target/internal/config"
	"speedtest/target/internal/service"
	"sync"
	"syscall"
)

func main() {
	cfg := config.Load()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup

	metrics := service.NewMetrics(cfg)
	metrics.Run(ctx, &wg)

	server := api.NewServer(cfg, metrics)
	server.Serve()

	wg.Wait()
}
