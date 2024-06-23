package cron

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron/v2"
	"sync"
	"time"
)

func Start(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	s, err := gocron.NewScheduler()
	if err != nil {
		panic(err)
	}

	// TODO: Implement the cron jobs
	j, err := s.NewJob(
		gocron.DurationJob(
			10*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("cron job")
			},
		),
	)
	if err != nil {
		panic(err)
	}

	fmt.Println(j.ID())

	s.Start()

	<-ctx.Done()

	err = s.Shutdown()
	if err != nil {
		panic(err)
	}

	fmt.Println("cron stopped")
}
