#!/bin/bash

# Ждём, пока Postgres точно запустится и начнёт слушать TCP
until pg_isready -h order_db -p 5435 -U postgres; do
  echo "Waiting for PostgreSQL to be ready..."
  sleep 1
done

# Проверка, есть ли уже подписка
psql -U postgres -d order -tc "SELECT 1 FROM pg_subscription WHERE subname = 'order_subscription'" | grep -q 1 || \
psql -U postgres -d order -c "CREATE SUBSCRIPTION order_subscription CONNECTION 'host=order_db port=5435 user=replicator password=replicator_password dbname=order' PUBLICATION menu_publication;"
