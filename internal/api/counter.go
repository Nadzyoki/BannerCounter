package api

import (
	"strings"
	"time"

	"github.com/Nadzyoki/BannerCounter/internal/utils"
	"github.com/valyala/fasthttp"
)

func handleCounter(counter Counter, ctx *fasthttp.RequestCtx) {
	id := strings.TrimPrefix(string(ctx.Path()), "/counter/")
	if id == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("missing id"))
		return
	}
	go counter.Add(utils.CreateKey(id, time.Now()))
	ctx.SetStatusCode(fasthttp.StatusOK)
}
