package saver

import (
	"context"
	"log/slog"
	"time"
)

type Storage interface {
	SaveCount(ctx context.Context, ids []string, counts []int) error
}

type Counter interface {
	GetAndReset() map[string]uint64
}

type Saver struct {
	counter  Counter
	storage  Storage
	interval time.Duration
}

func NewSaver(ctx context.Context, counter Counter, storage Storage, interval time.Duration) *Saver {
	saver := &Saver{
		counter:  counter,
		storage:  storage,
		interval: interval,
	}
	return saver
}

func (s *Saver) Schedule(ctx context.Context) {
	ticker := time.NewTicker(s.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.flush(ctx)
		case <-ctx.Done():
			s.flush(context.Background())
			return
		}
	}
}

func (s *Saver) flush(ctx context.Context) {
	mp := s.counter.GetAndReset()
	if len(mp) == 0 {
		return
	}

	ids := make([]string, 0, len(mp))
	counts := make([]int, 0, len(mp))
	for k, v := range mp {
		ids = append(ids, k)
		counts = append(counts, int(v))
	}

	if err := s.storage.SaveCount(ctx, ids, counts); err != nil {
		slog.Error("failed to save counts", "error", err)
	}
}
