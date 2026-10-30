package service

import (
	"context"
	"fmt"
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
		return fmt.Errorf("ping: invalid status code: %d", resp.StatusCode)
	}
	return nil
}
