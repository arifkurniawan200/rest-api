package app

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"net/http"
	"sync"
	"template/internal/usecase"
	"time"
)

type handler struct {
	User        usecase.UserUcase
	Transaction usecase.TransactionUcase
}

func Run(u usecase.UserUcase, t usecase.TransactionUcase) {
	e := echo.New()

	h := handler{
		User:        u,
		Transaction: t,
	}

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// CORS configuration
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept},
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	}))

	// Rate Limiter Configuration
	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 10, Burst: 30, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}
	e.Use(middleware.RateLimiterWithConfig(config))
	e.POST("/register", h.RegisterUser)
	e.POST("/login", h.LoginUser)

	customer := e.Group("/customer")
	{
		customer.Use(JWTMiddleware("secret")) // still default,can change anytime (i suggest i should placed in  .env)
	}
	admin := e.Group("/admin")
	{
		admin.Use(JWTMiddleware("secret")) // still default,can change anytime (i suggest i should placed in  .env)
		admin.Use(AdminMiddleware)
	}

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()

		if err := e.Start(":8080"); err != nil {
			fmt.Printf("Error starting server: %v\n", err)
		}
	}()

	wg.Wait()
}
