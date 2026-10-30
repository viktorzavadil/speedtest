package service

import (
	"context"
	"errors"
	"log"
	"math"
	"os"
	"speedtest/target/internal/config"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

type Metrics struct {
	cfg      *config.Config
	cpuUsage atomic.Uint64
}

type MetricsSnapshot struct {
	CpuUsage float64
}

func NewMetrics(cfg *config.Config) *Metrics {
	return &Metrics{cfg: cfg}
}

func (m *Metrics) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go m.runCpuTicker(ctx, wg)
}

func (m *Metrics) GetMetrics() MetricsSnapshot {
	return MetricsSnapshot{
		CpuUsage: math.Float64frombits(m.cpuUsage.Load()),
	}
}

func (m *Metrics) runCpuTicker(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	idle0, total0, err := readCpuTime()
	if err != nil {
		log.Printf("error reading cpu time: %v", err)
		return
	}
	t := time.NewTicker(m.cfg.ReadCpuInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			idle1, total1, err := readCpuTime()
			if err != nil || total1 == total0 {
				continue
			}
			u := 1 - float64(idle1-idle0)/float64(total1-total0)
			m.cpuUsage.Store(math.Float64bits(u))
			idle0, total0 = idle1, total1
		}
	}
}

func readCpuTime() (idle, total uint64, err error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) >= 9 && fields[0] == "cpu" {
			var sum uint64
			for _, f := range fields[1:9] {
				v, _ := strconv.ParseUint(f, 10, 64)
				sum += v
			}
			idleVal, _ := strconv.ParseUint(fields[4], 10, 64)
			iowaitVal, _ := strconv.ParseUint(fields[5], 10, 64)
			idleTotal := idleVal + iowaitVal
			return idleTotal, sum, nil
		}
	}
	return 0, 0, errors.New("cpu line not found")
}
