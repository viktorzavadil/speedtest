package service

import "time"

type Statistics struct {
	LastLatency   time.Duration
	AvgLatency    time.Duration
	PingTotal     int64
	PingSuccesses int64
	PingErrors    int64
	PingLost      float64

	TargetCpuUsage    float64
	TargetAllocMemory float64

	RxCount  int64
	RxErrors int64
	TxCount  int64
	TxErrors int64
}

type MetricsResponse struct {
	CPU struct {
		Usage float64 `json:"usage"`
	} `json:"cpu"`
	Memory struct {
		AllocMB    float64 `json:"allocMB"`
		Goroutines int     `json:"goroutines"`
	} `json:"memory"`
}
