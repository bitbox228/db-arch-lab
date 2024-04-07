package faker

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"sync"
)

const Count = 1_000_000
const WorkersCount = 100
const TablesCount = 1

func GenerateFakeData(pool *pgxpool.Pool) {
	var wg sync.WaitGroup

	wg.Add(1)
	go fillUsers(pool, &wg)

	wg.Wait()
}
