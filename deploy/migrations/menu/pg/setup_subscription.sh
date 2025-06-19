#!/bin/bash

# Ждём, пока Postgres точно запустится и начнёт слушать TCP
until pg_isready -h menu_db -p 5434 -U postgres; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 1
done

# Проверка, есть ли уже подписка
psql -U postgres -d menu -tc "SELECT 1 FROM pg_subscription WHERE subname = 'menu_subscription'" | grep -q 1 || \
psql -U postgres -d menu -c "CREATE SUBSCRIPTION menu_subscription CONNECTION 'host=menu_db port=5434 user=replicator password=replicator_password dbname=menu' PUBLICATION menu_publication;"
