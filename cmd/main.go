package main

import (
	"context"
	"db-arch-lab2/internal/fakers"
	"db-arch-lab2/internal/migrations"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

const DatabaseUrl = "postgres://postgres:postgres@db:5432/postgres?sslmode=disable"
const ConnectionsCount = 40

func main() {
	config, err := pgxpool.ParseConfig(DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	config.MaxConns = ConnectionsCount

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v\n", err)
	}
	if err := migrations.Migrate("db/migration", conn.Conn()); err != nil {
		log.Fatal(err)
	}
	conn.Release()
	log.Println("migrations done")

	fakers.GenerateFakeData(pool)
	log.Println("filled with fake data")
}
