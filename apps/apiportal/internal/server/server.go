package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/token"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	Port  int
	Store services.Store
	Config utils.Config
	Token token.TokenMaker
}

func NewServer(config utils.Config) *http.Server {
	ctx := context.Background()
	port, _ := strconv.Atoi(os.Getenv("PORT"))


	connPool, err := pgxpool.New(ctx, config.DB_URI)
	if err != nil {
		log.Fatal(err)
	}

	if err := connPool.Ping(ctx); err != nil {
		log.Fatal("cannot connect to db:", err)
	}

	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		log.Fatal(err)
	}

	newServer := &Server{
		Port:  port,
		Store: services.NewStore(connPool),
		Config: config,
		Token: tokenMaker,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.Port),
		Handler:      RegisterRoutes(newServer),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
