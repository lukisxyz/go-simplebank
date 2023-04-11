package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	db "github.com/flukis/simplebank/db/sqlc"
	"github.com/flukis/simplebank/util"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type Server struct {
	store           db.Store
	router          *echo.Echo
	validate        *validator.Validate
	passwordHashing util.Argon2Param
	tokenMaker      util.JWTMaker
	config          util.Config
}

func NewServer(store db.Store, cfg util.Config) (*Server, error) {
	v := validator.New()
	arg := util.Argon2Param{
		Memory:      64 * 1024,
		Iterations:  3,
		Parallelism: 2,
		SaltLength:  16,
		KeyLength:   32,
	}

	tokenMaker, err := util.NewJWTMaker(cfg.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		store:           store,
		validate:        v,
		passwordHashing: arg,
		tokenMaker:      *tokenMaker,
		config:          cfg,
	}
	serverRouter(server)
	return server, nil
}

func serverRouter(server *Server) {
	router := echo.New()

	router.POST("/account", server.CreateAccount)
	router.GET("/account/:id", server.GetAccount)
	router.GET("/account", server.FetchAccount)

	router.POST("/transfer", server.CreateTransfer)

	router.POST("/user", server.CreateUser)
	router.POST("/login", server.LoginUser)

	server.router = router
}

func (s *Server) Start(addr string) {
	go func() {
		if err := s.router.Start(addr); err != nil && err != http.ErrServerClosed {
			s.router.Logger.Fatal("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := s.router.Shutdown(ctx); err != nil {
		s.router.Logger.Fatal(err)
	}
}

type Meta struct {
	Limit int32 `json:"limit"`
	Page  int32 `json:"page"`
}
