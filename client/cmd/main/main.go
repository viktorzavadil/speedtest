package main

import (
	"context"
	_ "embed"
	"fmt"
	"log"
	"os/signal"
	"speedtest/client/internal/config"
	"speedtest/client/internal/service"
	"speedtest/client/internal/utils"
	"sync"
	"syscall"
	"time"
)

//go:embed banner.txt
var banner string

const (
	MB    = 1024 * 1024
	lines = 13
)

func main() {
	cfg := config.Load()

	fmt.Print("\033[34m\033[1m", banner, "\033[0m\n")

	fmt.Printf("\033[1m%20s:\033[0m %s\n", "target", cfg.TargetUrl)
	fmt.Printf("\033[1m%20s:\033[0m %d\n", "streams", cfg.Streams)
	fmt.Printf("\033[1m%20s:\033[0m %s\n", "idle phase duration", cfg.IdlePhaseDuration.String())
	fmt.Printf("\033[1m%20s:\033[0m %s\n", "rx phase duration", cfg.DownloadPhaseDuration.String())
	fmt.Printf("\033[1m%20s:\033[0m %s\n\n", "tx phase duration", cfg.UploadPhaseDuration.String())

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	duration := cfg.IdlePhaseDuration + cfg.DownloadPhaseDuration + cfg.UploadPhaseDuration + cfg.GracePeriod
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	var wg sync.WaitGroup

	target := service.NewTarget(cfg)

	err := target.Ping(ctx)
	if err != nil {
		log.Fatal("Failed to ping target:", err)
	}

	test := service.NewTest(cfg, target)

	wg.Add(1)
	go runStatistics(ctx, &wg, cfg, test)

	test.Run(ctx, &wg)

	wg.Wait()

	printFinalResult(cfg, test)
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
			return
		case <-ticker.C:
			stats := test.GetStatistics()

			now := time.Now()
			rx := stats.RxCount
			tx := stats.TxCount
			rxRate := utils.FormatBitrate(rx-prevRx, now.Sub(prevTime))
			txRate := utils.FormatBitrate(tx-prevTx, now.Sub(prevTime))
			prevRx, prevTx, prevTime = rx, tx, now

			var latencyStatistics service.LatencyStatistics
			if service.IsLoadedStatus(stats.TestStatus) {
				latencyStatistics = stats.LoadedLatency
			} else {
				latencyStatistics = stats.IdleLatency
			}

			if !firstRun {
				fmt.Printf("\033[%dA", lines)
			}
			firstRun = false

			fmt.Printf("\r\033[2K\033[1m%20s:\033[0m %s\n\n", "status", toStrStatus(stats.TestStatus))

			fmt.Printf("\r\033[2K%20s: %s\n", "last latency", utils.FormatDuration(latencyStatistics.LastLatency))
			fmt.Printf("\r\033[2K%20s: %s\n", "avg latency", utils.FormatDuration(latencyStatistics.AvgLatency))
			fmt.Printf("\r\033[2K%20s: %.1f%%\n", "ping lost", latencyStatistics.PingLost)
			fmt.Printf("\r\033[2K%20s: %s\n", "rx bytes", utils.FormatBytes(stats.RxCount))
			fmt.Printf("\r\033[2K%20s: %s\n", "rx bitrate", rxRate)
			fmt.Printf("\r\033[2K%20s: %d\n", "rx errors", stats.RxErrors)
			fmt.Printf("\r\033[2K%20s: %s\n", "tx bytes", utils.FormatBytes(stats.TxCount))
			fmt.Printf("\r\033[2K%20s: %s\n", "tx bitrate", txRate)
			fmt.Printf("\r\033[2K%20s: %d\n", "tx errors", stats.TxErrors)
			fmt.Printf("\r\033[2K%20s: %.1f%%\n", "target cpu usage", stats.TargetCpuUsage)
			fmt.Printf("\r\033[2K%20s: %.1fMB\n", "target mem usage", stats.TargetAllocMemory)
		}
	}
}

func printFinalResult(cfg *config.Config, test *service.Test) {
	stats := test.GetStatistics()

	fmt.Printf("\033[%dA\033[J", lines)

	fmt.Printf("\r\033[1m%20s:\033[0m %s\n\n", "status", toStrStatus(stats.TestStatus))

	fmt.Printf("%20s: min=%-10s avg=%-10s max=%s\n", "idle latency", utils.FormatDuration(stats.IdleLatency.MinLatency), utils.FormatDuration(stats.IdleLatency.AvgLatency), utils.FormatDuration(stats.IdleLatency.MaxLatency))
	fmt.Printf("%20s: min=%-10s avg=%-10s max=%s\n", "loaded latency", utils.FormatDuration(stats.LoadedLatency.MinLatency), utils.FormatDuration(stats.LoadedLatency.AvgLatency), utils.FormatDuration(stats.LoadedLatency.MaxLatency))
	fmt.Printf("%20s: %-10s %s\n", "rx", utils.FormatBytes(stats.RxCount), utils.FormatBitrate(stats.RxCount, cfg.DownloadPhaseDuration))
	fmt.Printf("%20s: %-10s %s\n", "tx", utils.FormatBytes(stats.TxCount), utils.FormatBitrate(stats.TxCount, cfg.UploadPhaseDuration))
}

func toStrStatus(status service.TestStatus) string {
	switch {
	case status == service.Initialization:
		return "initialization"
	case status == service.IdlePhase:
		return "idle phase"
	case status == service.DownloadPhase:
		return "download phase"
	case status == service.UploadPhase:
		return "upload phase"
	case status == service.Finished:
		return "finished"
	case status == service.Interrupted:
		return "interrupted"
	default:
		return "unknown"
	}
}
