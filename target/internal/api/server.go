package api

import (
	"fmt"
	"log"
	"net/http"
	"speedtest/target/internal/config"
	"speedtest/target/internal/service"
)

type Server struct {
	cfg     *config.Config
	metrics *service.Metrics

	mux *http.ServeMux
}

func NewServer(cfg *config.Config, metrics *service.Metrics) *Server {

	mux := http.NewServeMux()

	server := &Server{
		cfg:     cfg,
		metrics: metrics,
		mux:     mux,
	}

	mux.HandleFunc("GET /ping", ping)
	mux.HandleFunc("GET /metrics", server.getMetrics)
	mux.HandleFunc("GET /download", download)
	mux.HandleFunc("POST /upload", upload)

	return server
}

func (s *Server) Serve() {
	log.Println("server listening on", s.cfg.Port)
	err := http.ListenAndServe(fmt.Sprintf(":%d", s.cfg.Port), auth(s.cfg.ApiKey, s.mux))
	if err != nil {
		log.Fatal("initialization server failed", err)
	}
}
