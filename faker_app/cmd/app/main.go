package main

import (
	"context"
	"faker_app/internal/fakers"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
)

type DbConfig struct {
	Port     string `yaml:"port" env:"DB_PORT" env_default:"5432"`
	Host     string `yaml:"host" env:"DB_HOST" env_default:"localhost"`
	Name     string `yaml:"name" env:"POSTGRES_DB" env_default:"postgres"`
	User     string `yaml:"user" env:"POSTGRES_USER" env_default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env_default:"postgres"`
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

	if err := fakers.GenerateFakeData(pool); err != nil {
		log.Fatal(err)
	}
	log.Println("filled with fake data")

	pool.Close()
}
