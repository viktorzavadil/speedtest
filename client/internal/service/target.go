package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"speedtest/client/internal/config"
)

type Target struct {
	cfg *config.Config
}

func NewTarget(cfg *config.Config) *Target {
	return &Target{cfg: cfg}
}

func (t *Target) Ping(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", t.cfg.TargetUrl+"/ping", nil)
	if err != nil {
		return err
	}
	req.Header.Set("X-Api-Key", t.cfg.ApiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("GET ping: invalid status code: %d", resp.StatusCode)
	}
	return nil
}

func (t *Target) Metrics(ctx context.Context) (*MetricsResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", t.cfg.TargetUrl+"/metrics", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", t.cfg.ApiKey)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GET metrics: invalid status code: %d", resp.StatusCode)
	}
	var response MetricsResponse
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("GET metrics: decode response: %w", err)
	}
	return &response, nil
}

func (t *Target) Download(ctx context.Context) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", t.cfg.TargetUrl+"/download", nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", t.cfg.ApiKey)
	resp, err := http.DefaultClient.Do(req)
	return resp, err
}

func (t *Target) Upload(ctx context.Context, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, "POST", t.cfg.TargetUrl+"/upload", body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-Api-Key", t.cfg.ApiKey)
	req.Header.Set("Content-Type", "application/octet-stream")
	return http.DefaultClient.Do(req)
}
