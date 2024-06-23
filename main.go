package main

import (
	"context"
	"cyclic/cron"
	"cyclic/pkg/colonel"
	"cyclic/pkg/scribe"
	"cyclic/pkg/secretary"
	"cyclic/router"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func init() {
	colonel.Init()   // Initialize the configuration
	secretary.Init() // Initialize the database
	scribe.Init()    // Initialize the logger

	gin.SetMode(colonel.Writ.Server.Mode) // Set gin mode
}

func main() {
	ctx, cancel := context.WithCancel(context.Background()) // this context will be used to stop the server
	wg := &sync.WaitGroup{}                                 // wait group to wait for the server to stop gracefully

	go router.Route(ctx, wg) // Start the router
	go cron.Start(ctx, wg)   // Start the cron
	// TODO: Implement the message queue for mailer

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGQUIT)
	<-quit // Wait for the signal
	cancel()

	// Wait for the router and cron stop gracefully
	wg.Add(2) // add 2 because we started 2 goroutines above
	wg.Wait()

	fmt.Println("cyclic stopped") // TODO: Implement logging
	scribe.Scribe.Info("server stop gracefully")
	scribe.Scribe.Fatal("server stop gracefully")
}
