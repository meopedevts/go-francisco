package db

import (
	"context"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartDb() *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	poolConfig, err := pgxpool.ParseConfig("postgres://rinha:rinha2024@db:5432/rinha-db")
	if err != nil {
		log.Fatal(err.Error())
	}
	poolConfig.MaxConns = 128
	poolConfig.MinConns = 64

	conn, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatal(err.Error())
	}

	return conn
}
