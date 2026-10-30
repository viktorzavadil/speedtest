package api

import (
	"net/http"
	"runtime"
)

func metrics(w http.ResponseWriter, r *http.Request) {
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	resp := map[string]any{
		// todo
		//"cpu": map[string]float64{
		//	"current_pct": cpuCurrent,
		//	"avg_pct":     cpuAvg,
		//	"min_pct":     cpuMin,
		//	"max_pct":     cpuMax,
		//},
		"memory": map[string]any{
			"allocMB":    float64(mem.Alloc) / 1024 / 1024,
			"goroutines": runtime.NumGoroutine(),
		},
	}
	writeSuccess(w, resp)
}
