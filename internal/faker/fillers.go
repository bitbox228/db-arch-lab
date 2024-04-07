package faker

import (
	"context"
	"fmt"
	"github.com/go-faker/faker/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

func fillUsers(pool *pgxpool.Pool, wg *sync.WaitGroup) {
	defer wg.Done()
	userChan := make(chan user)
	for i := 0; i < WorkersCount/TablesCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for user := range userChan {
				pool.Exec(context.Background(),
					fmt.Sprintf(
						"INSERT INTO users (email, password_hash, nickname, avatar_url, is_private) VALUES ($1, $2, $3, $4, $5)"),
					user.Email,
					user.PasswordHash,
					user.Nickname,
					user.AvatarUrl,
					user.IsPrivate)
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
