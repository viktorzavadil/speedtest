package api

import (
	"net/http"
	"runtime"
)

func (s *Server) getMetrics(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	metrics := s.metrics.GetMetrics()

	resp := map[string]any{
		"cpu": map[string]float64{
			"usage": metrics.CpuUsage,
		},
		"memory": map[string]any{
			"allocMB":    float64(mem.Alloc) / 1024 / 1024,
			"goroutines": runtime.NumGoroutine(),
		},
	}
	writeSuccess(w, resp)
}
