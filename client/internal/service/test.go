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

	status atomic.Uint32

	lastLatency   atomic.Int64
	latencySum    atomic.Int64
	pingSuccesses atomic.Int64
	pingErrors    atomic.Int64
	latencyMu     sync.Mutex
	minLatency    time.Duration
	maxLatency    time.Duration

	targetCpuUsage    atomic.Uint64
	targetAllocMemory atomic.Uint64

	rxCount  atomic.Int64
	rxErrors atomic.Int64
	txCount  atomic.Int64
	txErrors atomic.Int64
}

func NewTest(cfg *config.Config, target *Target) *Test {
	test := &Test{cfg: cfg, target: target}
	test.status.Store(uint32(Initialization))
	return test
}

func (t *Test) Run(ctx context.Context, wg *sync.WaitGroup) {
	probeContext, probeCancel := context.WithCancel(ctx)

	wg.Add(1)
	go t.runLatencyProbe(probeContext, wg)
	wg.Add(1)
	go t.runTargetMetricsProbe(probeContext, wg)

	wg.Add(1)
	go func() {
		defer wg.Done()
		defer probeCancel()

		t.status.Store(uint32(DownloadPhase))
		t.runTestPhase(ctx, t.testDownload, t.cfg.DownloadPhaseDuration)
		t.status.Store(uint32(UploadPhase))
		t.runTestPhase(ctx, t.testUpload, t.cfg.UploadPhaseDuration)

		if ctx.Err() == nil {
			t.status.Store(uint32(Finished))
		} else {
			t.status.Store(uint32(Interrupted))
		}
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
		TestStatus: TestStatus(t.status.Load()),

		LastLatency:   time.Duration(t.lastLatency.Load()) * time.Microsecond,
		AvgLatency:    avgLatency,
		PingSuccesses: pingSuccesses,
		PingErrors:    pingErrors,
		PingTotal:     pingTotal,
		PingLost:      pingLost,

		MinLatency: t.minLatency,
		MaxLatency: t.maxLatency,

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
