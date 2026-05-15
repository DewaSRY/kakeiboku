package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	api_validator "github.com/dewasurya/kakeiboku/apps/apiportal/internal/validator"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/token"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/utils"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
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

	connPool, err := pgxpool.New(ctx, config.DB_URI)
	if err != nil {
		log.Fatal(err)
	}else{
		log.Println("Successfully connected to the database")
	}

	if err := connPool.Ping(ctx); err != nil {
		log.Fatal("cannot connect to db:", err)
	}else {
		log.Println("Successfully pinged the database")
	}

	tokenMaker, err := token.NewJWTMaker(config.SecretKey)
	if err != nil {
		log.Fatal(err)
	}

	newServer := &Server{
		Port:  config.Port,
		Store: services.NewStore(connPool),
		Config: config,
		Token: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", api_validator.ValidCurrency)
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.Port),
		Handler:      RegisterRoutes(newServer),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
