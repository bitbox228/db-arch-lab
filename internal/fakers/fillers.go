package fakers

import (
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
)

func fillUsers(pool *pgxpool.Pool, wg *sync.WaitGroup) {
	defer wg.Done()

	userChan := make(chan user)

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := &pgx.Batch{}

			for user := range userChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO users (email, password_hash, nickname, avatar_url, is_private) VALUES($1, $2, $3, $4, $5)"),
					user.Email,
					user.PasswordHash,
					user.Nickname,
					user.AvatarUrl,
					user.IsPrivate)
			}

			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var user user
	for i := 0; i < Count; i++ {
		faker.FakeData(&user)
		userChan <- user
	}

	close(userChan)
}

func fillAnime(pool *pgxpool.Pool, wg *sync.WaitGroup) {
	defer wg.Done()

	animeChan := make(chan anime)

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := &pgx.Batch{}

			for anime := range animeChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO anime (title, release_date, rating, genre, type, studio, status, age_rating, cover_url) VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9)"),
					anime.Title,
					anime.ReleaseDate,
					anime.Rating,
					anime.Genre,
					anime.Type,
					anime.Studio,
					anime.Status,
					anime.AgeRating,
					anime.CoverUrl)
			}

			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var anime anime
	for i := 0; i < Count; i++ {
		faker.FakeData(&anime)
		animeChan <- anime
	}

	close(animeChan)
}

func fillUserAnimeStatus(pool *pgxpool.Pool, wg *sync.WaitGroup) {
	defer wg.Done()

	userAnimeStatusChan := make(chan userAnimeStatus)

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := &pgx.Batch{}

			for userAnimeStatus := range userAnimeStatusChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO user_anime_status (anime_id, user_id, list, is_subscribed) VALUES($1, $2, $3, $4)"),
					userAnimeStatus.AnimeId,
					userAnimeStatus.UserId,
					userAnimeStatus.List,
					userAnimeStatus.IsSubscribed)
			}

			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var userAnimeStatus userAnimeStatus
	for i := 0; i < Count; i++ {
		faker.FakeData(&userAnimeStatus)
		userAnimeStatusChan <- userAnimeStatus
	}

	close(userAnimeStatusChan)
}
