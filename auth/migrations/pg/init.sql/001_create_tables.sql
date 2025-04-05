CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_login TIMESTAMPTZ
);

CREATE TABLE staff (
    id SERIAL PRIMARY KEY,
    work_email VARCHAR(255) UNIQUE NOT NULL, -- Отдельный email для работы
    work_phone VARCHAR(20) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL, -- Отдельный хеш пароля
    full_name VARCHAR(100) NOT NULL,
    hire_date DATE NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL, -- 'customer', 'admin', 'chef'
    description TEXT
);
