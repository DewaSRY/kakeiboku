package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/database"
	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port  int
	db    database.Service
	store services.Store
	config utils.Config
}

func NewServer() *http.Server {
	ctx := context.Background()
	config, err := utils.LoadConfig("../../")
	port, _ := strconv.Atoi(os.Getenv("PORT"))

	connPool, err := pgxpool.New(ctx, config.DB_URI)
	if err != nil {
		log.Fatal(err)
	}

	NewServer := &Server{
		port:  port,
		db:    database.New(),
		store: services.NewStore(connPool),
		config: config,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", config.Port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
