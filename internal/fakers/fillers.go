package fakers

import (
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
	"time"
)

func checkTable(pool *pgxpool.Pool, tableName string) bool {
	var count int
	err := pool.QueryRow(context.Background(), fmt.Sprintf("SELECT COUNT(*) FROM %s", tableName)).Scan(&count)
	if err != nil {
		return false
	}
	return count == 0
}

func fillUsers(pool *pgxpool.Pool, wg *sync.WaitGroup, innerWg *sync.WaitGroup) {
	defer wg.Done()
	if !checkTable(pool, "users") {
		return
	}

	userChan := make(chan user)

	for i := 0; i < WorkersCount; i++ {
		innerWg.Add(1)
		go func() {
			defer innerWg.Done()
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

func fillAnime(pool *pgxpool.Pool, wg *sync.WaitGroup, innerWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "anime") {
		return
	}

	animeChan := make(chan anime)

	for i := 0; i < WorkersCount; i++ {
		innerWg.Add(1)
		go func() {
			defer innerWg.Done()
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

func fillUserAnimeStatus(pool *pgxpool.Pool, wg *sync.WaitGroup, animeWg *sync.WaitGroup, userWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "user_anime_status") {
		return
	}

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

			userWg.Wait()
			animeWg.Wait()
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

func fillAnimeSeries(pool *pgxpool.Pool, wg *sync.WaitGroup, animeWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "anime_series") {
		return
	}

	animeSeriesChan := make(chan animeSeries)

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := &pgx.Batch{}

			for animeSeries := range animeSeriesChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO anime_series (anime_id, series_url, seconds_count) VALUES($1, $2, $3)"),
					animeSeries.AnimeId,
					animeSeries.SeriesUrl,
					animeSeries.SecondsCount)
			}

			animeWg.Wait()
			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var animeSeries animeSeries
	for i := 0; i < Count; i++ {
		faker.FakeData(&animeSeries)
		animeSeriesChan <- animeSeries
	}

	close(animeSeriesChan)
}

func fillReviews(pool *pgxpool.Pool, wg *sync.WaitGroup, innerWg *sync.WaitGroup, animeWg *sync.WaitGroup, userWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "reviews") {
		return
	}

	reviewsChan := make(chan reviews)

	for i := 0; i < WorkersCount; i++ {
		innerWg.Add(1)
		go func() {
			defer innerWg.Done()
			batch := &pgx.Batch{}

			for reviews := range reviewsChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO reviews (anime_id, user_id, rating, text) VALUES($1, $2, $3, $4)"),
					reviews.AnimeId,
					reviews.UserId,
					reviews.Rating,
					reviews.Text)
			}

			animeWg.Wait()
			userWg.Wait()
			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var reviews reviews
	for i := 0; i < Count; i++ {
		faker.FakeData(&reviews)
		reviewsChan <- reviews
	}

	close(reviewsChan)
}

func fillFriends(pool *pgxpool.Pool, wg *sync.WaitGroup, userWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "friends") {
		return
	}

	friendsChan := make(chan friends)

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := &pgx.Batch{}

			for friends := range friendsChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO friends (user_id1, user_id2) VALUES($1, $2)"),
					friends.UserId1,
					friends.UserId2)
			}

			userWg.Wait()
			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var friends friends
	for i := 0; i < Count; i++ {
		faker.FakeData(&friends)
		friendsChan <- friends
	}

	close(friendsChan)
}

func fillMessages(pool *pgxpool.Pool, wg *sync.WaitGroup, userWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "messages") {
		return
	}

	messagesChan := make(chan messages)

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := &pgx.Batch{}

			for messages := range messagesChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO messages (sender_id, receiver_id, text, file_url, time) VALUES($1, $2, $3, $4, $5)"),
					messages.SenderId,
					messages.ReceiverId,
					messages.Text,
					messages.FileUrl,
					time.Now())
			}

			userWg.Wait()
			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var messages messages
	for i := 0; i < Count; i++ {
		faker.FakeData(&messages)
		messagesChan <- messages
	}

	close(messagesChan)
}

func fillAchievements(pool *pgxpool.Pool, wg *sync.WaitGroup, innerWg *sync.WaitGroup, animeWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "achievements") {
		return
	}

	achievementsChan := make(chan achievements)

	for i := 0; i < WorkersCount; i++ {
		innerWg.Add(1)
		go func() {
			defer innerWg.Done()
			batch := &pgx.Batch{}

			for achievements := range achievementsChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO achievements (anime_id, name, description) VALUES($1, $2, $3)"),
					achievements.AnimeId,
					achievements.Name,
					achievements.Description)
			}

			animeWg.Wait()
			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var achievements achievements
	for i := 0; i < Count; i++ {
		faker.FakeData(&achievements)
		achievementsChan <- achievements
	}

	close(achievementsChan)
}

func fillUserAchievements(pool *pgxpool.Pool, wg *sync.WaitGroup, achievementsWg *sync.WaitGroup, userWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "user_achievements") {
		return
	}

	userAchievementsChan := make(chan userAchievements)

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := &pgx.Batch{}

			for userAchievements := range userAchievementsChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO user_achievements (achievement_id, user_id, time) VALUES($1, $2, $3)"),
					userAchievements.AchievementId,
					userAchievements.UserId,
					time.Now())
			}

			userWg.Wait()
			achievementsWg.Wait()
			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var userAchievements userAchievements
	for i := 0; i < Count; i++ {
		faker.FakeData(&userAchievements)
		userAchievementsChan <- userAchievements
	}

	close(userAchievementsChan)
}

func fillNotifications(pool *pgxpool.Pool, wg *sync.WaitGroup, userWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "notifications") {
		return
	}

	notificationsChan := make(chan notifications)

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := &pgx.Batch{}

			for notifications := range notificationsChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO notifications (user_id, type, body, time) VALUES($1, $2, $3, $4)"),
					notifications.UserId,
					notifications.Type,
					"{}",
					time.Now())
			}

			userWg.Wait()
			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var notifications notifications
	for i := 0; i < Count; i++ {
		faker.FakeData(&notifications)
		notificationsChan <- notifications
	}

	close(notificationsChan)
}

func fillReactions(pool *pgxpool.Pool, wg *sync.WaitGroup, reviewWg *sync.WaitGroup, userWg *sync.WaitGroup) {
	defer wg.Done()

	if !checkTable(pool, "reactions") {
		return
	}

	reactionsChan := make(chan reactions)

	for i := 0; i < WorkersCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := &pgx.Batch{}

			for reactions := range reactionsChan {
				batch.Queue(fmt.Sprintf(
					"INSERT INTO reactions (review_id, user_id, is_like) VALUES($1, $2, $3)"),
					reactions.ReviewId,
					reactions.UserId,
					reactions.IsLike)
			}

			reviewWg.Wait()
			userWg.Wait()
			br := pool.SendBatch(context.Background(), batch)

			if err := br.Close(); err != nil {
				log.Fatal(err)
			}
		}()
	}

	var reactions reactions
	for i := 0; i < Count; i++ {
		faker.FakeData(&reactions)
		reactionsChan <- reactions
	}

	close(reactionsChan)
}
