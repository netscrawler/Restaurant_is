-- Создание пользователя для репликации, если он не существует
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'replicator') THEN
      CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD 'replicator_password';
   END IF;
END $$;

-- Настройка прав доступа
GRANT CONNECT ON DATABASE menu TO replicator;
GRANT USAGE ON SCHEMA public TO replicator;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO replicator;

-- Создание публикации, если она не существует
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_publication WHERE pubname = 'menu_publication') THEN
      CREATE PUBLICATION menu_publication FOR ALL TABLES;
   END IF;
END $$;
