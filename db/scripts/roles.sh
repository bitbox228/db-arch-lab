#!/bin/bash

psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "CREATE ROLE reader LOGIN;"
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "CREATE ROLE writer LOGIN;"
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "CREATE ROLE group_role NOLOGIN;"
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "CREATE USER analytic WITH PASSWORD 'analytic';"

psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "GRANT SELECT ON TABLE $ANALYTICS_TABLE to analytic;"
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "GRANT SELECT ON ALL TABLES IN SCHEMA public TO reader;"
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "GRANT SELECT, INSERT, UPDATE ON ALL TABLES IN SCHEMA public TO writer;"
psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO group_role;"

IFS=','
users=($USERS)

for user in "${users[@]}"; do
  psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "CREATE USER $user WITH PASSWORD '$user';"
  psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c "GRANT group_role to $user;"
done

