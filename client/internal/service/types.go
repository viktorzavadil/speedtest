package service

type MetricsResponse struct {
	CPU struct {
		Usage float64 `json:"usage"`
	} `json:"cpu"`
	Memory struct {
		AllocMB    float64 `json:"allocMB"`
		Goroutines int     `json:"goroutines"`
	} `json:"memory"`
}
