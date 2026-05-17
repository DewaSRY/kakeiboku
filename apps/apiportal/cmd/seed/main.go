package main

import (
	"context"

	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/services"
	"github.com/dewasurya/kakeiboku/apps/apiportal/pkg/utils"
	"github.com/dewasurya/kakeiboku/apps/apiportal/seed"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	config, err := utils.LoadConfig(".")
	if err != nil {
		panic("cannot load config")
	}

	conn_pool, err := pgxpool.New(ctx, config.DB_URI)

	if err != nil {
		panic("cannot create connection pool")
	}

	store := services.NewStore(conn_pool)

	seed.Seed001(ctx, store)

}
