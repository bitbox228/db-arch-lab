package main

import (
	"context"
	"db-arch-lab2/internal/migrations"
	"github.com/jackc/pgx/v5"
	"log"
)

func main() {
	dbUrl := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"
	db, err := pgx.Connect(context.Background(), dbUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close(context.Background())

	if err := migrations.Migrate("db/migration", db); err != nil {
		log.Fatal(err)
	}
}
