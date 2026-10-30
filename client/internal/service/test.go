package service

import (
	"context"
	"math"
	"speedtest/client/internal/config"
	"sync"
	"sync/atomic"
	"time"
)

type Test struct {
	cfg    *config.Config
	target *Target

	lastLatency   atomic.Int64
	latencySum    atomic.Int64
	pingSuccesses atomic.Int64
	pingErrors    atomic.Int64
	latencyMu     sync.Mutex
	latencies     []time.Duration

	targetCpuUsage    atomic.Uint64
	targetAllocMemory atomic.Uint64
}

type Statistics struct {
	LastLatency   time.Duration
	AvgLatency    time.Duration
	PingTotal     int64
	PingSuccesses int64
	PingErrors    int64
	PingLost      float64

	TargetCpuUsage    float64
	TargetAllocMemory float64
}

func NewTest(cfg *config.Config, target *Target) *Test {
	return &Test{cfg: cfg, target: target}
}

func (t *Test) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go t.runLatencyProbe(ctx, wg)
	wg.Add(1)
	go t.runTargetMetricsProbe(ctx, wg)
}

func (t *Test) GetStatistics() Statistics {
	pingSuccesses := t.pingSuccesses.Load()
	pingErrors := t.pingErrors.Load()
	pingTotal := pingSuccesses + pingErrors
	var pingLost float64 = 0
	if pingTotal != 0 {
		pingLost = float64(pingErrors) / float64(pingTotal) * 100
	}
	latencySum := t.latencySum.Load()
	var avgLatency time.Duration
	if pingSuccesses > 0 {
		avgLatency = time.Duration(latencySum/pingSuccesses) * time.Microsecond
	}
	return Statistics{
		LastLatency:   time.Duration(t.lastLatency.Load()) * time.Microsecond,
		AvgLatency:    avgLatency,
		PingSuccesses: pingSuccesses,
		PingErrors:    pingErrors,
		PingTotal:     pingTotal,
		PingLost:      pingLost,

		TargetCpuUsage:    math.Float64frombits(t.targetCpuUsage.Load()) * 100,
		TargetAllocMemory: math.Float64frombits(t.targetAllocMemory.Load()),
	}
}

func (t *Test) runLatencyProbe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(t.cfg.LatencyProbeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			start := time.Now()
			err := t.target.Ping(ctx)
			if err != nil {
				if ctx.Err() == nil {
					t.pingErrors.Add(1)
				}
				continue
			}
			t.pingSuccesses.Add(1)
			latency := time.Since(start)
			t.lastLatency.Store(latency.Microseconds())
			t.latencySum.Add(latency.Microseconds())
			t.latencyMu.Lock()
			t.latencies = append(t.latencies, latency)
			t.latencyMu.Unlock()
		}
	}
}

func (t *Test) runTargetMetricsProbe(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	ticker := time.NewTicker(t.cfg.TargetMetricsProbeInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			response, err := t.target.Metrics(ctx)
			if err != nil {
				continue
			}
			t.targetCpuUsage.Store(math.Float64bits(response.CPU.Usage))
			t.targetCpuUsage.Store(math.Float64bits(response.CPU.Usage))
			t.targetAllocMemory.Store(math.Float64bits(response.Memory.AllocMB))
		}
	}
}
