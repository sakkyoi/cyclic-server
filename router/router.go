package router

import (
	"context"
	"cyclic/middleware/jwt"
	"cyclic/pkg/colonel"
	"cyclic/router/api"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"sync"
)

func Route(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	server := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", colonel.Writ.Server.Host, colonel.Writ.Server.Port),
		Handler: build(),
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(err)
		}
	}()

	<-ctx.Done() // wait for the context to be done

	fmt.Println("router stopped")
}

func build() *gin.Engine {
	r := gin.New()
	if err := r.SetTrustedProxies(colonel.Writ.Server.TrustedProxies); err != nil {
		panic(err)
	}
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	a := api.New()

	apiGroup := r.Group("/api")
	{
		apiGroup.GET("/ding", a.Ding)

		apiGroup.GET("/signup", a.CheckIsSignupEnabled)    // check if signup is enabled
		apiGroup.POST("/signup", a.Signup)                 // signup
		apiGroup.PUT("/signup", jwt.JWT(), a.VerifySignup) // verify signup with token if verification is enabled

		apiGroup.POST("/login", a.Login)

		apiGroup.GET("/user", jwt.JWT(), a.User.Get)
	}

	return r
}
