package app

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"template/config"
	"template/internal/usecase"
	"time"
)

type handler struct {
	User        usecase.UserUcase
	Transaction usecase.TransactionUcase
}

var cv = NewCustomValidator()

func Run(u usecase.UserUcase, t usecase.TransactionUcase) {
	e := echo.New()

	h := handler{
		User:        u,
		Transaction: t,
	}

	cfg := config.ReadConfig()

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
	rateLimiterConfig := middleware.RateLimiterConfig{
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

	v1 := e.Group("/v1")
	v1.Use(middleware.RateLimiterWithConfig(rateLimiterConfig))
	v1.POST("/register", h.RegisterUser)
	v1.POST("/login", h.LoginUser)

	// set versioning v1

	user := v1.Group("/user")
	{
		user.Use(JWTMiddleware(cfg.Env.SecretKey))
	}

	item := v1.Group("/items")
	{
		item.Use(JWTMiddleware(cfg.Env.SecretKey))
	}

	admin := v1.Group("/admin")
	{
		admin.Use(JWTMiddleware(cfg.Env.SecretKey))
		admin.Use(AdminMiddleware)
	}

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	go func() {

		if err := e.Start(":8080"); err != nil && err != http.ErrServerClosed {
			log.Printf("Error starting server: %v\n", err)
			os.Exit(1)
		}
	}()

	<-stop

	log.Println("OS SIGNAL RECEIVED")

	ctx, cancel := context.WithTimeout(context.Background(), 40*time.Second)

	defer cancel()

	log.Println("SHUTTING DOWN SERVER...")

	if err := e.Shutdown(ctx); err != nil {
		log.Printf("ERR SHUTTING DOWN SERVER : %v\n", err)
		os.Exit(1)
	}

	log.Println("SERVER GRACEFULLY STOPPED")
}
