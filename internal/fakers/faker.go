package fakers

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"strconv"
	"sync"
)

var Count = 1_000_000

const WorkersCount = 10

func GenerateFakeData(pool *pgxpool.Pool) error {
	count, err := strconv.Atoi(os.Getenv("COUNT"))
	if err != nil {
		return err
	}
	Count = count
	log.Printf("%d lines in every table", Count)

	addIdFaker()
	var wg, userWg, animeWg, reviewWg, achievementWg sync.WaitGroup

	go fillUsers(pool, &wg, &userWg)
	go fillAnime(pool, &wg, &animeWg)
	go fillAnimeSeries(pool, &wg, &animeWg)
	wg.Wait()
	log.Println(1)

	go fillAchievements(pool, &wg, &achievementWg, &animeWg)
	go fillReviews(pool, &wg, &reviewWg, &animeWg, &userWg)
	go fillUserAnimeStatus(pool, &wg, &userWg, &animeWg)
	go fillFriends(pool, &wg, &userWg)
	wg.Wait()
	log.Println(2)

	go fillMessages(pool, &wg, &userWg)
	go fillNotifications(pool, &wg, &userWg)
	go fillUserAchievements(pool, &wg, &achievementWg, &userWg)
	go fillReactions(pool, &wg, &reviewWg, &userWg)
	wg.Wait()
	log.Println(3)

	return nil
}
