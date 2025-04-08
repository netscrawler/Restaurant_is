---------------------------
-- 1. Таблица пользователей (клиенты)
---------------------------
CREATE TABLE clients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    phone VARCHAR(255) UNIQUE,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
);
--TODO: Добавить номер телефона в таблицы
---------------------------
-- 2. Таблица сотрудников (ресторан)
---------------------------
CREATE TABLE staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ
);

---------------------------
-- 3. OAuth-провайдеры (для клиентов)
---------------------------
CREATE TABLE oauth_providers (
    client_id UUID NOT NULL REFERENCES clients(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL CHECK (provider IN ('yandex', 'google')),
    provider_id VARCHAR(255) NOT NULL,
    access_token TEXT NOT NULL,
    PRIMARY KEY (client_id, provider)
);

---------------------------
-- 4. Refresh-токены (общая таблица)
---------------------------
CREATE TABLE refresh_tokens (
    token TEXT PRIMARY KEY,
    user_id UUID NOT NULL, -- Может ссылаться на clients.id или staff.id
    user_type VARCHAR(10) NOT NULL CHECK (user_type IN ('client', 'staff')),
    expires_at TIMESTAMPTZ NOT NULL,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

---------------------------
-- 5. Аудит авторизации
---------------------------
CREATE TABLE auth_logs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,         -- ID из таблицы clients/staff
    user_type VARCHAR(10) NOT NULL CHECK (user_type IN ('client', 'staff')),
    action VARCHAR(20) NOT NULL CHECK (action IN ('login', 'logout', 'token_refresh', 'token_revoke')),
    ip_address VARCHAR(45) NOT NULL,
    user_agent TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
