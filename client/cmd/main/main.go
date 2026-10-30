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

const MB = 1024 * 1024

func main() {
	cfg := config.Load()

	fmt.Print("\033[34m\033[1m", banner, "\033[0m\n")

	fmt.Printf("\033[1m            target:\033[0m %s\n", cfg.TargetUrl)
	fmt.Printf("\033[1m           streams:\033[0m %d parallel streams\n", cfg.Streams)
	fmt.Printf("\033[1m rx phase duration:\033[0m %s\n", cfg.DownloadPhaseDuration.String())
	fmt.Printf("\033[1m tx phase duration:\033[0m %s\n\n", cfg.UploadPhaseDuration.String())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	duration := cfg.IdlePhaseDuration + cfg.DownloadPhaseDuration + cfg.UploadPhaseDuration
	ctx, cancel := context.WithTimeout(ctx, duration)
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

	prevTime := time.Now()
	var (
		prevRx int64
		prevTx int64
	)

	for {
		select {
		case <-ctx.Done():
			fmt.Println()
			return
		case <-ticker.C:
			stats := test.GetStatistics()

			now := time.Now()
			rx := stats.RxCount
			tx := stats.TxCount
			rxRate := toMbit(float64(rx-prevRx)) / now.Sub(prevTime).Seconds()
			txRate := toMbit(float64(tx-prevTx)) / now.Sub(prevTime).Seconds()
			prevRx, prevTx, prevTime = rx, tx, now

			if !firstRun {
				fmt.Printf("\033[%dA", 11)
			}
			firstRun = false
			fmt.Printf("\r\033[2K      last latency: %s\n", stats.LastLatency)
			fmt.Printf("\r\033[2K       avg latency: %s\n", stats.AvgLatency)
			fmt.Printf("\r\033[2K         ping lost: %.1f%%\n", stats.PingLost)
			fmt.Printf("\r\033[2K          rx bytes: %.1fMB\n", float64(stats.RxCount)/MB)
			fmt.Printf("\r\033[2K           rx rate: %.1fMbit/s\n", rxRate)
			fmt.Printf("\r\033[2K         rx errors: %d\n", stats.RxErrors)
			fmt.Printf("\r\033[2K          tx bytes: %.1fMB\n", float64(stats.TxCount)/MB)
			fmt.Printf("\r\033[2K           tx rate: %.1fMbit/s\n", txRate)
			fmt.Printf("\r\033[2K         tx errors: %d\n", stats.TxErrors)
			fmt.Printf("\r\033[2K  target cpu usage: %.1f%%\n", stats.TargetCpuUsage)
			fmt.Printf("\r\033[2K  target mem usage: %.1fMB\n", stats.TargetAllocMemory)
		}
	}
}

func toMbit(bytes float64) float64 {
	return bytes * 8 / 1e6
}
