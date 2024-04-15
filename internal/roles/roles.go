package roles

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log"
	"os"
	"strings"
)

func AddRoles(ctx context.Context, db *pgx.Conn) error {
	_, err := db.Exec(ctx, `
        CREATE ROLE reader LOGIN;
        CREATE ROLE writer LOGIN;
        CREATE ROLE group_role NOLOGIN;
    `)
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, fmt.Sprintf(`
		CREATE USER analytic WITH PASSWORD 'analytic';
		GRANT SELECT ON TABLE %s to analytic;
	`, os.Getenv("ANALYTICS_TABLE")))
	if err != nil {
		return err
	}

	_, err = db.Exec(ctx, `
        GRANT SELECT ON ALL TABLES IN SCHEMA public TO reader;
        GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO writer;
        GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO group_role;
    `)
	if err != nil {
		return err
	}

	groupUsers := strings.Split(os.Getenv("USERS"), ",")
	log.Println(groupUsers)
	for _, user := range groupUsers {
		_, err = db.Exec(ctx, fmt.Sprintf(`
			CREATE USER %s WITH PASSWORD '%s';
			GRANT group_role to %s;
		`, user, user, user))
		if err != nil {
			return err
		}
	}

	return nil
}
