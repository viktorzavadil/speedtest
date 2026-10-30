package service

import (
	"context"
	"io"
	"net/http"
	"speedtest/client/internal/utils"
)

func (t *Test) testUpload(ctx context.Context) {
	resp, err := t.target.Upload(ctx, utils.NewZeroReader(ctx, &t.txCount))
	if err != nil {
		if ctx.Err() == nil {
			t.txErrors.Add(1)
		}
		return
	}
	defer resp.Body.Close()
	io.Copy(io.Discard, resp.Body)
	if resp.StatusCode != http.StatusOK && ctx.Err() == nil {
		t.txErrors.Add(1)
	}
}
