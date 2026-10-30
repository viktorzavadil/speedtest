package service

import "time"

type TestStatus uint32

const (
	Initialization TestStatus = iota
	IdlePhase      TestStatus = iota
	DownloadPhase  TestStatus = iota
	UploadPhase    TestStatus = iota
	Finished       TestStatus = iota
	Interrupted    TestStatus = iota
)

type LatencyStatistics struct {
	LastLatency   time.Duration
	AvgLatency    time.Duration
	MinLatency    time.Duration
	MaxLatency    time.Duration
	PingTotal     int64
	PingSuccesses int64
	PingErrors    int64
	PingLost      float64
}

type Statistics struct {
	TestStatus TestStatus

	IdleLatency   LatencyStatistics
	LoadedLatency LatencyStatistics

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
