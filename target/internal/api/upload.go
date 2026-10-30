package api

import (
	"io"
	"net/http"
)

func upload(w http.ResponseWriter, r *http.Request) {
	n, err := io.Copy(io.Discard, r.Body)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.Header().Set("Content-Type", "application/json")
	writeSuccess(w, map[string]int64{"receivedBytes": n})
}
