package service

import "time"

type TestStatus uint32

const (
	Initialization TestStatus = iota
	DownloadPhase  TestStatus = iota
	UploadPhase    TestStatus = iota
	Finished       TestStatus = iota
	Interrupted    TestStatus = iota
)

type Statistics struct {
	TestStatus TestStatus

	LastLatency   time.Duration
	AvgLatency    time.Duration
	PingTotal     int64
	PingSuccesses int64
	PingErrors    int64
	PingLost      float64

	MinLatency time.Duration
	MaxLatency time.Duration

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
