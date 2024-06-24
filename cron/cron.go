package cron

import (
	"context"
	"cyclic/pkg/scribe"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"github.com/samber/lo"
	"go.uber.org/zap"
	"sync"
	"time"
)

func Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done() // send signal to the wait group that this goroutine is done

	s, err := gocron.NewScheduler()
	if err != nil {
		scribe.Scribe.Fatal("failed to create cron scheduler", zap.Error(err))
	}

	// TODO: Implement the cron jobs
	_, err = s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				scribe.Scribe.Debug("cron job triggered")
			},
		),
	)
	if err != nil {
		scribe.Scribe.Fatal("failed to create cron job", zap.Error(err))
	}

	s.Start()

	scribe.Scribe.Info("cron started")
	scribe.Scribe.Debug(fmt.Sprintf("started %d cron job", len(s.Jobs())), zap.Strings("jobs", lo.Map(s.Jobs(), func(job gocron.Job, index int) string { return job.ID().String() })))

	<-ctx.Done()

	err = s.Shutdown()
	if err != nil {
		scribe.Scribe.Error("failed to shutdown cron scheduler", zap.Error(err)) // use error instead of fatal because we want to continue the other goroutines to stop gracefully
	}

	scribe.Scribe.Info("cron stopped")
}
