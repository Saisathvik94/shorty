package workers

import (
	"context"
	"log"
	"time"

	"github.com/Saisathvik94/shorty/apps/api/internal/repository"
)

type CleanUpWorker struct {
	repo     *repository.URLRepository
	interval time.Duration
}

func NewCleanUpWorker(repo *repository.URLRepository, interval time.Duration) *CleanUpWorker {
	return &CleanUpWorker{
		repo:     repo,
		interval: interval,
	}
}

func (w *CleanUpWorker) Start(ctx context.Context) {
	ticker := time.NewTicker(w.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("cleanup worker stopped!")
			return
		case <-ticker.C:
			count, err := w.repo.DeactivateExpiredURLs(ctx)

			if err != nil {
				log.Printf("cleanup worker failed:%v", err)
				continue
			}

			log.Printf("cleanup worker deactivated %d expired URLs", count)
		}
	}
}
