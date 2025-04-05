---------------------------
-- 1. Таблица типов аккаунтов
---------------------------
CREATE TABLE account_types (
    id SMALLSERIAL PRIMARY KEY,
    type_name VARCHAR(50) UNIQUE NOT NULL
);

---------------------------
-- 2. Таблица ролей
---------------------------
CREATE TABLE roles (
    id SMALLSERIAL PRIMARY KEY,
    role_name VARCHAR(50) UNIQUE NOT NULL
);

---------------------------
-- 3. Таблица пользователей
---------------------------
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    account_type_id SMALLINT NOT NULL REFERENCES account_types(id),
    email VARCHAR(255) UNIQUE,
    work_email VARCHAR(255) UNIQUE,
    phone VARCHAR(20) UNIQUE,
    password_hash VARCHAR(255),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ,

    -- Проверки для разных типов аккаунтов
    CONSTRAINT email_check CHECK (
        (account_type_id = 1 AND email IS NOT NULL) OR  -- USER
        (account_type_id = 2 AND work_email IS NOT NULL) -- STAFF
    )
);

-- Индексы
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_work_email ON users(work_email);

---------------------------
-- 4. Связь пользователей с ролями
---------------------------
CREATE TABLE user_roles (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id SMALLINT NOT NULL REFERENCES roles(id),
    PRIMARY KEY (user_id, role_id)
);

---------------------------
-- 5. OAuth-провайдеры
---------------------------
CREATE TABLE oauth_providers (
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL, -- yandex, google и т.д.
    provider_user_id VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL,
    refresh_token TEXT,
    expires_at TIMESTAMPTZ,
    PRIMARY KEY (user_id, provider)
);

---------------------------
-- 6. Refresh-токены
---------------------------
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

---------------------------
-- 7. Данные курьеров (только для STAFF)
---------------------------
CREATE TABLE courier_details (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    vehicle_type VARCHAR(50),      -- Тип транспорта: car, bike, scooter
    vehicle_number VARCHAR(20),    -- Номер транспорта
    work_zone VARCHAR(255)         -- Зона работы
);
