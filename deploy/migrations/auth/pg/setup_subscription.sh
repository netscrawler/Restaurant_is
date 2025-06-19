#!/bin/bash

# Ждём, пока Postgres точно запустится и начнёт слушать TCP
until pg_isready -h auth_db -p 5432 -U postgres; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 1
done

# Проверка, есть ли уже подписка
psql -U postgres -d auth -tc "SELECT 1 FROM pg_subscription WHERE subname = 'auth_subscription'" | grep -q 1 || \
psql -U postgres -d auth -c "CREATE SUBSCRIPTION auth_subscription CONNECTION 'host=auth_db port=5432 user=replicator password=replicator_password dbname=auth' PUBLICATION auth_publication;"
