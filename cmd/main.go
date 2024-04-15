package main

import (
	"context"
	"db-arch-lab2/internal/fakers"
	"db-arch-lab2/internal/migrations"
	"db-arch-lab2/internal/roles"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DbConfig struct {
	Port     string `yaml:"port" env:"PORT" env_default:"5432"`
	Host     string `yaml:"host" env:"HOST" env_default:"localhost"`
	Name     string `yaml:"name" env:"DB_NAME" env_default:"postgres"`
	User     string `yaml:"user" env:"DB_USER" env_default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env_default:"postgres"`
	SslMode  string `yaml:"sslMode" env:"SSL_MODE" env_default:"disable"`
}

const ConnectionsCount = 40

func main() {
	var cfg DbConfig

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(err)
	}

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SslMode)

	config, err := pgxpool.ParseConfig(databaseUrl)
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

	if err := fakers.GenerateFakeData(pool); err != nil {
		log.Fatal(err)
	}
	log.Println("filled with fake data")

	conn, err = pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Unable to acquire a database connection: %v\n", err)
	}

	if err := roles.AddRoles(context.Background(), conn.Conn()); err != nil {
		log.Fatal(err)
	}
	conn.Release()
	log.Println("added roles")

	pool.Close()
}
