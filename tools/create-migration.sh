#!/bin/bash

if [ -z "$1" ]; then 
    echo "Error: Please provide a migration name"
    echo "Usage: $0 <migration-name>"
    exit 1
fi

mkdir -p db/migrations

timestamp=$(date '+%Y%m%d%H%M%S')

migration_name=$(echo "$1" | tr '[:upper:]' '[:lower:]' | tr ' ' '_')

up_file="db/migrations/${timestamp}_${migration_name}.up.sql"
down_file="db/migrations/${timestamp}_${migration_name}.down.sql"

touch "$up_file"
touch "$down_file"

echo "Created migration files: "
echo " $up_file"
echo " $down_file"
