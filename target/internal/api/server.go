package api

import (
	"fmt"
	"log"
	"net/http"
	"speedtest/target/internal/config"
)

var zeroBuffer = make([]byte, 32*1024)

type Server struct {
	cfg *config.Config
	mux *http.ServeMux
}

func NewServer(cfg *config.Config) *Server {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /ping", ping)
	mux.HandleFunc("GET /metrics", metrics)
	mux.HandleFunc("GET /download", download)
	mux.HandleFunc("POST /upload", upload)

	return &Server{
		cfg: cfg,
		mux: mux,
	}
}

func (s *Server) Serve() {
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), auth(s.cfg.ApiKey, s.mux))
	if err != nil {
		log.Fatal("initialization server failed", err)
	}
	log.Println("server listening on", s.cfg.Port)
}
