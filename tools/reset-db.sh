#!/bin/bash

if [ -z "$1" ]; then 
    echo "Usage: $0 <path-to-env-file>"
    echo "Example: $0 ./.env"
    exit 1
fi

ENV_FILE=$1

if [ ! -f "$ENV_FILE" ]; then 
    echo "Error: .env file was not found"
    exit 1
fi

# Load the environment variables
set -a
source "$ENV_FILE"
set +a

required_vars=("PG_HOST" "PG_PORT" "PG_USERNAME" "PG_PASSWORD" "PG_DBNAME")
for var in "${required_vars[@]}"; do 
    if [ -z "${!var}" ]; then
        echo "Error: $var is not set in the $ENV_FILE"
        exit 1
    fi
done
#
echo "Starting to reset database ${PG_DBNAME}..."

export PGPASSWORD=$PG_PASSWORD

if ! psql -h "$PG_HOST" \
     -p "$PG_PORT" \
     -U "$PG_USERNAME" \
     -d "postgres" \
     -c "DROP DATABASE IF EXISTS ${PG_DBNAME};"; then
     echo "Error: failed to drop the database"
     unset PGPASSWORD
     exit 1
fi

psql -h "$PG_HOST" \
     -p "$PG_PORT" \
     -U "$PG_USERNAME" \
     -d "postgres" \
     -c "CREATE DATABASE ${PG_DBNAME};"

unset PGPASSWORD

echo "Database ${PG_DBNAME} has be reset"

