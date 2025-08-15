package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Nadzyoki/BannerCounter/internal/api"
	"github.com/Nadzyoki/BannerCounter/internal/atomiccounter"
	"github.com/Nadzyoki/BannerCounter/internal/config"
	"github.com/Nadzyoki/BannerCounter/internal/logger"
	"github.com/Nadzyoki/BannerCounter/internal/repo"
	"github.com/Nadzyoki/BannerCounter/internal/saver"
	"github.com/valyala/fasthttp"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger.InitLogger()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	cfg, err := config.NewConfig()
	if err != nil {
		slog.Error(err.Error())
		return
	}
	storage, err := repo.NewRepo(ctx, cfg)
	if err != nil {
		slog.Error(err.Error())
		return
	}
	defer storage.Close()

	counter := atomiccounter.NewAtomicCounter()
	saver := saver.NewSaver(ctx, counter, storage, cfg.Interval)

	go saver.Schedule(ctx)

	server := &fasthttp.Server{
		Handler: api.NewMux(counter, storage),
	}

	go func() {
		if err := server.ListenAndServe(":" + cfg.ListenPort); err != nil {
			panic(err)
		}
	}()
	slog.Info("server started", "port", cfg.ListenPort)

	<-sigCh
	cancel()

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if err := server.ShutdownWithContext(shutdownCtx); err != nil {
		panic(err)
	}
}
