package main

import (
	"context"
	"cyclic/cron"
	"cyclic/mailer"
	"cyclic/pkg/colonel"
	"cyclic/pkg/dispatcher"
	"cyclic/pkg/magistrate"
	"cyclic/pkg/scribe"
	"cyclic/pkg/secretary"
	"cyclic/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func init() {
	// Initialize the configuration
	if err := colonel.Init(); err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize configuration: %v", err)) // cause the logger(scribe) is not initialized yet, we use log.Fatal instead
	}

	// Initialize the logger, logger must be initialized before anything else except the configuration
	if err := scribe.Init(); err != nil {
		log.Fatal(fmt.Sprintf("failed to initialize logger: %v", err)) // cause the logger(scribe) initialization failed, we use log.Fatal instead
	}

	// Initialize the database
	if err := secretary.Init(); err != nil {
		scribe.Scribe.Fatal("failed to initialize database", zap.Error(err))
	}

	// Initialize the message queue
	if err := dispatcher.Init(); err != nil {
		scribe.Scribe.Fatal("failed to initialize message queue", zap.Error(err))
	}

	gin.SetMode(colonel.Writ.Server.Mode) // Set gin mode

	// test the magistrate(to avoid invalid keys)
	_ = magistrate.New()
}

func main() {
	scribe.Scribe.Info("server start")
	ctx, cancel := context.WithCancel(context.Background()) // this context will be used to stop the server
	wg := &sync.WaitGroup{}                                 // wait group to wait for the server to stop gracefully

	go router.Route(ctx, wg) // Start the router
	go cron.Start(ctx, wg)   // Start the cron
	go mailer.Start(ctx, wg) // Start the mailer

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
	<-quit // Wait for the signal
	cancel()

	// Wait for the router and cron stop gracefully
	wg.Add(3) // add 3 because we started 3 goroutines above
	wg.Wait()

	scribe.Scribe.Info("server stop gracefully")
}
