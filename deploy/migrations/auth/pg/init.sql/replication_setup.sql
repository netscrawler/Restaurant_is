-- Создание пользователя для репликации
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = 'replicator') THEN
      CREATE ROLE replicator WITH REPLICATION LOGIN PASSWORD 'replicator_password';
   END IF;
END $$;

-- Настройка прав
GRANT CONNECT ON DATABASE auth TO replicator;
GRANT USAGE ON SCHEMA public TO replicator;
GRANT SELECT ON ALL TABLES IN SCHEMA public TO replicator;
ALTER DEFAULT PRIVILEGES IN SCHEMA public GRANT SELECT ON TABLES TO replicator;

-- Создание публикации
DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_publication WHERE pubname = 'auth_publication') THEN
      CREATE PUBLICATION auth_publication FOR ALL TABLES;
   END IF;
END $$;

