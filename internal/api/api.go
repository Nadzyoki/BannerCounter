package api

import (
	"context"
	"strings"

	"github.com/valyala/fasthttp"
)

type Counter interface {
	Add(id string)
	Get(id string) uint64
}

type Storage interface {
	GetStats(ctx context.Context, bannerID string, tsFrom, tsTo string) (map[string]int, error)
}

func NewMux(counter Counter, storage Storage) func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		path := string(ctx.Path())

		switch {
		case strings.HasPrefix(path, "/counter/"):
			if !ctx.IsGet() {
				ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
				return
			}
			handleCounter(counter, ctx)
		case strings.HasPrefix(path, "/stats/"):
			if !ctx.IsPost() {
				ctx.SetStatusCode(fasthttp.StatusMethodNotAllowed)
				return
			}
			handleStats(storage, ctx)
		default:
			ctx.SetStatusCode(fasthttp.StatusNotFound)
		}
	}
}
