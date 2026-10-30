package api

import (
	"crypto/subtle"
	"net/http"
)

func auth(apiKey string, next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		clientApiKey := r.Header.Get("X-Api-Key")
		if subtle.ConstantTimeCompare([]byte(clientApiKey), []byte(apiKey)) != 1 {
			writeError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		next.ServeHTTP(w, r)
	}
}
