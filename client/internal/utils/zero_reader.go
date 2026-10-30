package utils

import (
	"context"
	"sync/atomic"
)

type ZeroReader struct {
	ctx     context.Context
	counter *atomic.Int64
}

func NewZeroReader(ctx context.Context, counter *atomic.Int64) *ZeroReader {
	return &ZeroReader{ctx: ctx, counter: counter}
}

func (r *ZeroReader) Read(p []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}
	clear(p)
	r.counter.Add(int64(len(p)))
	return len(p), nil
}
