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

	rxCount  atomic.Int64
	rxErrors atomic.Int64
	txCount  atomic.Int64
	txErrors atomic.Int64
}

func NewTest(cfg *config.Config, target *Target) *Test {
	return &Test{cfg: cfg, target: target}
}

func (t *Test) Run(ctx context.Context, wg *sync.WaitGroup) {
	wg.Add(1)
	go t.runLatencyProbe(ctx, wg)
	wg.Add(1)
	go t.runTargetMetricsProbe(ctx, wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		t.runTestPhase(ctx, t.testDownload, t.cfg.DownloadPhaseDuration)
		t.runTestPhase(ctx, t.testUpload, t.cfg.UploadPhaseDuration)
	}()
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

		RxCount:  t.rxCount.Load(),
		RxErrors: t.rxErrors.Load(),
		TxCount:  t.txCount.Load(),
		TxErrors: t.txErrors.Load(),
	}
}

func (t *Test) runTestPhase(ctx context.Context, phase func(ctx context.Context), duration time.Duration) {
	ctx, cancel := context.WithTimeout(ctx, duration)
	defer cancel()

	var phaseWg sync.WaitGroup

	for i := 0; i < int(t.cfg.Streams); i++ {
		phaseWg.Add(1)
		go func() {
			defer phaseWg.Done()
			phase(ctx)
		}()
	}
	phaseWg.Wait()
}
