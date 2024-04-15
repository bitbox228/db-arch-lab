package migrations

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type migration struct {
	Version  string
	FileName string
}

const versionEnv = "MIGRATION_VERSION"

var migrationPattern = regexp.MustCompile(`^V(\d+\.\d+\.\d+)_(.+)\.sql$`)

func Migrate(directory string, db *pgx.Conn) error {
	files, err := os.ReadDir(directory)
	if err != nil {
		log.Fatal(err)
	}

	var migrations []migration

	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".sql" {
			matches := migrationPattern.FindStringSubmatch(file.Name())
			if len(matches) == 3 {
				migrations = append(migrations, migration{
					Version:  matches[1],
					FileName: file.Name(),
				})
			}
		}
	}

	sort.Slice(migrations, func(i, j int) bool {
		return compareVersions(migrations[i].Version, migrations[j].Version)
	})

	version := os.Getenv(versionEnv)

	for _, migration := range migrations {
		if len(version) != 0 && compareVersions(version, migration.Version) {
			break
		}

		query, err := os.ReadFile(directory + "/" + migration.FileName)
		if err != nil {
			return err
		}

		_, err = db.Exec(context.Background(), string(query))
		if err != nil {
			return err
		}
	}

	return nil
}

func compareVersions(version1, version2 string) bool {
	version1Splitted := strings.Split(version1, ".")
	version2Splitted := strings.Split(version2, ".")
	for k := 0; k < 3; k++ {
		if version1Splitted[k] != version2Splitted[k] {
			v1, _ := strconv.Atoi(version1Splitted[k])
			v2, _ := strconv.Atoi(version2Splitted[k])
			return v1 < v2
		}
	}
	return false
}
