package service

import (
	"context"
	"math"
	"sync"
	"time"
)

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
			t.targetAllocMemory.Store(math.Float64bits(response.Memory.AllocMB))
		}
	}
}
