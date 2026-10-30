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
			status := TestStatus(t.status.Load())
			var latency *Latency
			if IsLoadedStatus(status) {
				latency = t.loadedLatency
			} else {
				latency = t.idleLatency
			}
			start := time.Now()
			err := t.target.Ping(ctx)
			if err != nil {
				if ctx.Err() == nil {
					latency.pingErrors.Add(1)
				}
				continue
			}
			latency.pingSuccesses.Add(1)
			latencyDur := time.Since(start)
			latency.lastLatency.Store(latencyDur.Microseconds())
			latency.latencySum.Add(latencyDur.Microseconds())
			latency.latencyMu.Lock()
			if latency.minLatency == 0 || latencyDur < latency.minLatency {
				latency.minLatency = latencyDur
			}
			if latencyDur > latency.maxLatency {
				latency.maxLatency = latencyDur
			}
			latency.latencyMu.Unlock()
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
