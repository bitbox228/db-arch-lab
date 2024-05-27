package main

import (
	"context"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"math/rand/v2"
	"os"
	"regexp"
	"slices"
	"strconv"
	"time"
)

var queries = []string{`WITH anime_ranking_by_genre AS (SELECT a.genre,
                                       a.id                                                         as anime_id,
                                       a.title                                                      as anime_title,
                                       COUNT(DISTINCT ual.user_id)                                  AS views_count,
                                       ROW_NUMBER()
                                       OVER (PARTITION BY a.genre ORDER BY COUNT(ual.user_id) DESC) AS genre_rank
                                FROM anime a
                                         JOIN user_anime_status ual
                                              ON a.id = ual.anime_id AND ual.list <> 'WANT_TO_WATCH'
                                GROUP BY a.genre, a.id
                                ORDER BY a.genre, genre_rank), anime_reviews_stats AS (SELECT a.id AS id,
                                    COUNT(DISTINCT r.id)                        AS reviews_count,
                                    SUM(CASE WHEN rr.is_like THEN 1 ELSE 0 END) AS reviews_likes_count
                             FROM anime a
                                      LEFT JOIN
                                  reviews r ON a.id = r.anime_id
                                      LEFT JOIN reactions rr
                                                ON r.id = rr.review_id
                             GROUP BY a.id)

SELECT rg.genre, rg.anime_id, rg.anime_title, rg.views_count, ars.reviews_count, ars.reviews_likes_count
FROM anime_ranking_by_genre rg
         JOIN anime_reviews_stats ars
              ON rg.anime_id = ars.id
WHERE rg.genre_rank <= 10
ORDER BY rg.genre, rg.genre_rank;`,

	`WITH RECURSIVE friends_of_friends AS (SELECT f.user_id2 as user_id, 0 AS level
                                      FROM friends f
                                      WHERE user_id1 = $1
                                      UNION ALL
                                      SELECT f.user_id2 AS user_id, fof.level + 1 AS level
                                      FROM friends f
                                               JOIN friends_of_friends fof ON fof.user_id = f.user_id1
                                      WHERE fof.level < 1)
SELECT DISTINCT fof.user_id,
                fof.level
FROM friends_of_friends fof
ORDER BY fof.level, fof.user_id;`,

	`SELECT DATE_TRUNC('month', ua.time) AS month_year,
       COUNT(*)                     AS achievements_count
FROM users u
         JOIN user_achievements ua ON u.id = ua.user_id
         JOIN achievements a ON a.id = ua.achievement_id
WHERE ua.time >= NOW() - INTERVAL '5 years'
GROUP BY DATE_TRUNC('month', ua.time)
ORDER BY month_year;`,
	`WITH anime_ranking_by_genre AS (SELECT a.genre,
                                       a.id                                                         as anime_id,
                                       a.title                                                      as anime_title,
                                       COUNT(DISTINCT ual.user_id)                                  AS views_count,
                                       ROW_NUMBER()
                                       OVER (PARTITION BY a.genre ORDER BY COUNT(ual.user_id) DESC) AS genre_rank
                                FROM anime a
                                         JOIN user_anime_status ual
                                              ON a.id = ual.anime_id AND ual.list <> 'WANT_TO_WATCH'
                                GROUP BY a.genre, a.id
                                ORDER BY a.genre, genre_rank),
     anime_reviews_stats AS (SELECT a.id                                        AS id,
                                    COUNT(DISTINCT r.id)                        AS reviews_count,
                                    SUM(CASE WHEN rr.is_like THEN 1 ELSE 0 END) AS reviews_likes_count
                             FROM anime a
                                      LEFT JOIN
                                  reviews r ON a.id = r.anime_id
                                      LEFT JOIN reactions rr
                                                ON r.id = rr.review_id
                             GROUP BY a.id),
     reviews_stats AS (SELECT r.id                                        AS id,
                              r.anime_id                                  AS anime_id,
                              SUM(CASE WHEN rr.is_like THEN 1 ELSE 0 END) AS likes_count
                       FROM reviews r
                                LEFT JOIN reactions rr ON r.id = rr.review_id
                       GROUP BY r.id)

SELECT rg.genre,
       rg.genre_rank,
       rg.anime_id,
       rg.anime_title,
       rg.views_count,
       ars.reviews_count,
       ars.reviews_likes_count,
       rs.id          AS review_id,
       rs.likes_count AS likes_count
FROM anime_ranking_by_genre rg
         JOIN anime_reviews_stats ars
              ON rg.anime_id = ars.id
         LEFT JOIN reviews_stats rs
                   ON rg.anime_id = rs.anime_id
WHERE rg.genre_rank <= 10
ORDER BY rg.genre, rg.genre_rank;`}

const (
	pattern = `cost=\d+\.\d+\.{2}(\d+\.\d+)`
)

type DbConfig struct {
	Port     string `yaml:"port" env:"DB_PORT" env_default:"5432"`
	Host     string `yaml:"host" env:"DB_HOST" env_default:"localhost"`
	Name     string `yaml:"name" env:"POSTGRES_DB" env_default:"postgres"`
	User     string `yaml:"user" env:"POSTGRES_USER" env_default:"postgres"`
	Password string `yaml:"password" env:"POSTGRES_PASSWORD" env_default:"postgres"`
	SslMode  string `yaml:"sslMode" env:"SSL_MODE" env_default:"disable"`
}

type ExplainAnalyzeResult struct {
	Best  float64
	Worst float64
	Avg   float64
}

func main() {
	var cfg DbConfig

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatal(err)
	}

	queriesCount, err := strconv.Atoi(os.Getenv("ANALYZER_COUNT"))
	if err != nil {
		log.Fatal(err)
	}
	count, err := strconv.Atoi(os.Getenv("COUNT"))
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Create(fmt.Sprintf("logs/%s", time.Now().UTC().Format(time.DateTime)))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	databaseUrl := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SslMode)

	config, err := pgxpool.ParseConfig(databaseUrl)
	if err != nil {
		log.Fatal(err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatal(err)
	}
	defer pool.Close()
	re := regexp.MustCompile(pattern)

	results := make([]ExplainAnalyzeResult, 0, len(queries))

	for i, query := range queries {
		explainQuery := "EXPLAIN ANALYZE " + query
		costs := make([]float64, 0, queriesCount)
		for j := 0; j < queriesCount; j++ {
			var rows pgx.Rows
			if i == 1 {
				id := rand.IntN(count + 1)
				rows, err = pool.Query(context.Background(), explainQuery, id)
			} else {
				rows, err = pool.Query(context.Background(), explainQuery)
			}
			if err != nil {
				log.Fatal(err)
			}
			if rows.Next() {
				var s string
				if err = rows.Scan(&s); err != nil {
					log.Fatal(err)
				}
				cost, err := strconv.ParseFloat(re.FindStringSubmatch(s)[1], 64)
				if err != nil {
					log.Fatal(err)
				}
				costs = append(costs, cost)
			}
			rows.Close()
		}
		var sum float64
		for _, cost := range costs {
			sum += cost
		}
		log.Printf("query %d done", i)
		results = append(results, ExplainAnalyzeResult{
			Best:  slices.Min(costs),
			Worst: slices.Max(costs),
			Avg:   sum / float64(queriesCount),
		})
	}

	for i, result := range results {
		file.WriteString(fmt.Sprintf("%d\n", i))
		file.WriteString(fmt.Sprintf("best:  %f\n", result.Best))
		file.WriteString(fmt.Sprintf("worst: %f\n", result.Worst))
		file.WriteString(fmt.Sprintf("avg:   %f\n", result.Avg))
	}
}
