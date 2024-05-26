#!/bin/bash

echo "started backup process"

backup_file="backups/$(date +'%Y%m%d%H%M%S').sql"

pg_dump -h "${DB_HOST}" -U "${POSTGRES_USER}" "${POSTGRES_DB}" > "${backup_file}"

while [ $(ls -1 backups | wc -l) -gt ${BACKUPS_COUNT} ]; do
  oldest=$(ls -1 backups | head -n 1)
  rm backups/"${oldest}"
done

echo "finished backup process"