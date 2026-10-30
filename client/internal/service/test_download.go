package service

import (
	"context"
	"io"
	"net/http"
)

func (t *Test) testDownload(ctx context.Context) {
	resp, err := t.target.Download(ctx)
	if err != nil {
		if ctx.Err() == nil {
			t.rxErrors.Add(1)
		}
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.rxErrors.Add(1)
		return
	}
	buf := make([]byte, 64*1024)
	for {
		n, err := resp.Body.Read(buf)
		t.rxCount.Add(int64(n))
		if err != nil {
			if err != io.EOF && ctx.Err() == nil {
				t.rxErrors.Add(1)
			}
			return
		}
	}
}
