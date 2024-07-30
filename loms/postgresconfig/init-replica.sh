#!/bin/bash
set -e

export PGPASSWORD="$POSTGRES_PASSWORD"

echo "Очистка директории данных..."
rm -rf /var/lib/postgresql/data/*

until psql -h $POSTGRES_MASTER_HOST -p 5432 -U "$POSTGRES_USER" -d "$POSTGRES_DB" -c '\q'; do
  >&2 echo "Основной сервер не доступен - ожидание..."
  sleep 1
done

echo "Основной сервер доступен, инициализация реплики..."

pg_basebackup -h $POSTGRES_MASTER_HOST -p 5432 -U $POSTGRES_USER -D /var/lib/postgresql/data -Fp -Xs -P -R

echo "Репликация инициализирована"
