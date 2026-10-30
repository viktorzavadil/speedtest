package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os/signal"
	"speedtest/client/internal/config"
	"speedtest/client/internal/service"
	"sync"
	"syscall"
	"time"
)

//go:embed banner.txt
var banner string

func main() {
	cfg := config.Load()

	fmt.Print("\033[34m\033[1m", banner, "\033[0m\n")

	fmt.Printf("\033[1m           target:\033[0m %s\n", cfg.TargetUrl)
	fmt.Printf("\033[1m          streams:\033[0m %d parallel streams\n", cfg.Streams)
	fmt.Printf("\033[1m         duration:\033[0m %s per phase\n\n", cfg.Duration.String())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithTimeout(ctx, cfg.Duration)
	defer cancel()

	var wg sync.WaitGroup

	target := service.NewTarget(cfg)

	err := target.Ping(ctx)
	if err != nil {
		log.Fatal("Failed to ping target:", err)
	}

	test := service.NewTest(cfg, target)

	test.Run(ctx, &wg)

	wg.Add(1)
	go runStatistics(ctx, &wg, cfg, test)

	wg.Wait()
}

func runStatistics(ctx context.Context, wg *sync.WaitGroup, cfg *config.Config, test *service.Test) {
	defer wg.Done()

	ticker := time.NewTicker(cfg.StatisticsInterval)
	defer ticker.Stop()
	firstRun := true

	for {
		select {
		case <-ctx.Done():
			fmt.Println()
			return
		case <-ticker.C:
			stats := test.GetStatistics()
			if !firstRun {
				fmt.Printf("\033[%dA", 5)
			}
			firstRun = false
			fmt.Printf("\r\033[2K     last latency: %s\n", stats.LastLatency)
			fmt.Printf("\r\033[2K      avg latency: %s\n", stats.AvgLatency)
			fmt.Printf("\r\033[2K        ping lost: %.1f%%\n", stats.PingLost)
			fmt.Printf("\r\033[2K target cpu usage: %.1f%%\n", stats.TargetCpuUsage)
			fmt.Printf("\r\033[2K target mem usage: %.1fMB\n", stats.TargetAllocMemory)
		}
	}
}
