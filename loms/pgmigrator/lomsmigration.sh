#!/bin/bash
#source .env

export MIGRATION_DSN="host=$DB_HOST port=$DB_PORT dbname=$DB_NAME user=$DB_USER password=$DB_PASSWORD sslmode=disable"

sleep 2 && goose -dir "./migrations" postgres "${MIGRATION_DSN}" up -v
