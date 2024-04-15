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

	wg.Add(1)
	go fillUsers(pool, &wg, &userWg)
	wg.Add(1)
	go fillAnime(pool, &wg, &animeWg)
	wg.Add(1)
	go fillAnimeSeries(pool, &wg, &animeWg)
	wg.Wait()
	log.Println(1)

	wg.Add(1)
	go fillAchievements(pool, &wg, &achievementWg, &animeWg)
	wg.Add(1)
	go fillReviews(pool, &wg, &reviewWg, &animeWg, &userWg)
	wg.Add(1)
	go fillUserAnimeStatus(pool, &wg, &userWg, &animeWg)
	wg.Add(1)
	go fillFriends(pool, &wg, &userWg)
	wg.Wait()
	log.Println(2)

	wg.Add(1)
	go fillMessages(pool, &wg, &userWg)
	wg.Add(1)
	go fillNotifications(pool, &wg, &userWg)
	wg.Add(1)
	go fillUserAchievements(pool, &wg, &achievementWg, &userWg)
	wg.Add(1)
	go fillReactions(pool, &wg, &reviewWg, &userWg)
	wg.Wait()
	log.Println(3)

	return nil
}
