#!/bin/bash

# Ждём, пока Postgres точно запустится и начнёт слушать TCP
until pg_isready -h user_db -p 5439 -U postgres; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 1
done

# Проверка, есть ли уже подписка
psql -U postgres -d user -tc "SELECT 1 FROM pg_subscription WHERE subname = 'user_subscription'" | grep -q 1 || \
psql -U postgres -d user -c "CREATE SUBSCRIPTION order_subscription CONNECTION 'host=user_db port=5439 user=replicator password=replicator_password dbname=user' PUBLICATION user_publication;"
