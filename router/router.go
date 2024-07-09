package router

import (
	"context"
	"cyclic/middleware/jwt"
	"cyclic/pkg/colonel"
	"cyclic/pkg/scribe"
	"cyclic/router/api"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"sync"
)

func Route(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done() // send signal to the wait group that this goroutine is done

	router := build() // build the router

	// create the server
	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", colonel.Writ.Server.Host, colonel.Writ.Server.Port),
		Handler: router,
	}

	// start the server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			scribe.Scribe.Fatal("failed to start router", zap.Error(err))
		}
	}()

	scribe.Scribe.Info("router started")
	scribe.Scribe.Debug("router listening on", zap.String("address", fmt.Sprintf("%s:%d", colonel.Writ.Server.Host, colonel.Writ.Server.Port)))

	<-ctx.Done() // wait for the context to be done

	scribe.Scribe.Info("router stopped")
}

func build() *gin.Engine {
	r := gin.New()
	if err := r.SetTrustedProxies(colonel.Writ.Server.TrustedProxies); err != nil {
		scribe.Scribe.Fatal("failed to set trusted proxies", zap.Error(err))
	}

	scribe.Scribe.Debug("router trusted proxies", zap.Strings("proxies", colonel.Writ.Server.TrustedProxies))

	r.Use(gin.Logger())   // gin middleware to log the request into gin.DefaultWriter (it's an os.Stdout)
	r.Use(gin.Recovery()) // gin middleware to recover from any panic and return 500

	a := api.New()

	apiGroup := r.Group("/api")
	{
		// health check
		apiGroup.GET("/ding", a.Ding)

		// signup
		apiGroup.GET("/signup", a.Signup.Check)         // check if signup is enabled
		apiGroup.POST("/signup", a.Signup.Signup)       // signup
		apiGroup.PATCH("/signup", a.Signup.Resend)      // resend verification email
		apiGroup.GET("/verify/:token", a.Signup.Verify) // verify with token if verification is enabled

		// user
		apiGroup.POST("/auth", a.User.Auth)

		apiGroup.GET("/user", jwt.JWT(), a.User.Get)

		// plan
		apiGroup.POST("/plan", jwt.JWT(), a.Plan.Create)
		apiGroup.PUT("/plan/:id", jwt.JWT(), a.Plan.Update)
		apiGroup.DELETE("/plan/:id", jwt.JWT(), a.Plan.Remove)
		apiGroup.GET("/plan", jwt.JWT(), a.Plan.List)
		apiGroup.POST("/plan/:id/invite", jwt.JWT(), a.Plan.Invite)          // invite user to subscription
		apiGroup.GET("/plan/:id/subscribers", jwt.JWT(), a.Plan.Subscribers) // list subscribers of a subscription

		// subscription
		// TODO: list subscription (for user)
		// TODO: subscribe (if user is invited)
		// TODO: leave subscription (must not be the owner)
	}

	return r
}
