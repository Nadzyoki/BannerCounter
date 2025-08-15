package api

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/Nadzyoki/BannerCounter/models"
	"github.com/valyala/fasthttp"
)

func handleStats(storage Storage, ctx *fasthttp.RequestCtx) {
	id := strings.TrimPrefix(string(ctx.Path()), "/stats/")
	if id == "" {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("missing id"))
		return
	}
	var req models.StatsRequest
	if err := json.Unmarshal(ctx.PostBody(), &req); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		ctx.SetBody([]byte("invalid JSON body"))
		return
	}
	fromTime, err := time.Parse(time.RFC3339, req.From)
	if err != nil {
		fromTime, err = time.Parse("2006-01-02T15:04:05", req.From)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBody([]byte("invalid from date"))
			return
		}
	}

	toTime, err := time.Parse(time.RFC3339, req.To)
	if err != nil {
		toTime, err = time.Parse("2006-01-02T15:04:05", req.To)
		if err != nil {
			ctx.SetStatusCode(fasthttp.StatusBadRequest)
			ctx.SetBody([]byte("invalid to date"))
			return
		}
	}

	stats, err := storage.GetStats(ctx, id, fromTime.Format(time.RFC3339), toTime.Format(time.RFC3339))
	if err != nil {
		slog.Error("failed to get stats", "error", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("failed to get stats"))
		return
	}

	statsResp := make([]models.Stat, 0, len(stats))
	for tm, count := range stats {
		statsResp = append(statsResp, models.Stat{
			Ts: tm,
			V:  count,
		})
	}

	resp := map[string]interface{}{"stats": statsResp}
	body, err := json.Marshal(resp)
	if err != nil {
		slog.Error("failed to marshal response", "error", err)
		ctx.SetStatusCode(fasthttp.StatusInternalServerError)
		ctx.SetBody([]byte("failed to marshal response"))
		return
	}

	ctx.SetContentType("application/json")
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetBody(body)
}
