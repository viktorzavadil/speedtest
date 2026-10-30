package service

import (
	"context"
	"errors"
	"math"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

// todo

type CPUMonitor struct {
	usage atomic.Uint64 // float64 bity, aby šlo číst bez zámku
}

func (m *CPUMonitor) Run(ctx context.Context, interval time.Duration) {
	idle0, total0, _ := ReadCPUTicks()
	t := time.NewTicker(interval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			idle1, total1, err := ReadCPUTicks()
			if err != nil || total1 == total0 {
				continue
			}
			u := 1 - float64(idle1-idle0)/float64(total1-total0)
			m.usage.Store(math.Float64bits(u))
			idle0, total0 = idle1, total1
		}
	}
}

func (m *CPUMonitor) Usage() float64 {
	return math.Float64frombits(m.usage.Load())
}

func ReadCPUTicks() (idle, total uint64, err error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, 0, err
	}
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		fields := strings.Fields(line)
		if len(fields) > 4 && fields[0] == "cpu" {
			var sum uint64
			for _, f := range fields[1:9] {
				v, _ := strconv.ParseUint(f, 10, 64)
				sum += v
			}
			idleVal, _ := strconv.ParseUint(fields[4], 10, 64)
			return idleVal, sum, nil
		}
	}
	return 0, 0, errors.New("cpu line not found")
}
