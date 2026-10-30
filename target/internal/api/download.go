package api

import "net/http"

var zeroBuffer = make([]byte, 32*1024)

func download(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/octet-stream")
	for {
		select {
		case <-r.Context().Done():
			return
		default:
			if _, err := w.Write(zeroBuffer); err != nil {
				return
			}
		}
	}
}
