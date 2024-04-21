#!/bin/bash

function compare_versions() {
  local version1=$1
  local version2=$2

  local IFS='.'
  local version1_splitted=($version1)
  local version2_splitted=($version2)

  for ((k = 0; k < 3; k++)); do
    if [ "${version1_splitted[k]}" -lt "${version2_splitted[k]}" ]; then
      echo -1
      return
    elif [ "${version1_splitted[k]}" -gt "${version2_splitted[k]}" ]; then
      echo 1
      return
    fi
  done

  echo 0
  return
}

last_version=$MIGRATION_VERSION
regex='^V([0-9]+\.[0-9]+\.[0-9]+)_.+\.sql$'

for migration in $(ls db/migration/V*.sql | sort -V); do
  migration_file=$(basename "$migration")
  if [[ "$migration_file" =~ $regex ]]; then
    if [ "$last_version" = "" ]; then
      psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$migration"
      continue
    fi
    version="${BASH_REMATCH[1]}"
    result=$(compare_versions "$version" "$last_version")
    if [ "$result" -eq 1 ]; then
      break
    fi
    psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -f "$migration"
  fi
done