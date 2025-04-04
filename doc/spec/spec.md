= SPEC-1: Информационная система управления Доставкой
:sectnums:
:toc:

== Background

Современные сервисы доставки требуют эффективного управления заказами, оплатами, меню и пользователями. Информационная система должна обеспечивать взаимодействие клиентов, официантов, поваров и администраторов через удобный интерфейс. Важными аспектами являются мониторинг, логирование и аналитика данных.

== Requirements

=== Must-have (Обязательные)

- Регистрация и аутентификация пользователей (клиенты, администраторы, официанты, повара).
- Управление заказами (создание, обновление, статусы, уведомления).
- Интеграция с платежной системой для онлайн-оплат (заглушка).
- Управление меню (CRUD-операции с блюдами, категориями, ценами).
- Система отчетности и аналитики (реплики баз данных для BI-инструмента Metabase).
- Логирование событий и мониторинг работы системы.

=== Should-have (Желательные)

- Поддержка программ лояльности и скидок.
- Интеграция с внешними сервисами доставки.

=== Could-have (Дополнительные)

- Автоматизация обработки заказов с помощью ML-модели (например, предсказание популярности блюд). (Не важно, будет алгоритмом, без ML)

=== Won't-have (Не планируется на текущий этап)

- Поддержка нескольких ресторанов в одной системе.
- Голосовое управление заказами.

== Method

=== Архитектурный стиль

- Микросервисная архитектура.
- API Gateway для управления запросами.
- Асинхронное взаимодействие (Kafka / RabbitMQ для событий).
- Репликация БД для аналитики (Metabase).
- Паттерн SAGA для контроля за операциями

=== Технологический стек

- **Бэкенд:** Go (Gin, gRPC, pgxPool)
- **Фронтенд:** React (Next.js / Vite)
- **База данных:** PostgreSQL (отдельная БД на сервис), S3 для хранения фото блюд. Мб кликхаус на нагруженные бд вместо PostgreSQL
- **Сообщения:** Kafka
- **Мониторинг:** Prometheus + Grafana
- **Логирование:** Loki + Grafana
- **Развертывание** Docker + IaC(Ansible, Terraform)

=== Основные сервисы

0. **Notify** - Сервис для отправки уведомлений пользователям(будет предствалять из себя заглушку с тг ботом)
1. **AuthService** – Аутентификация (JWT). Принимает логин, пароль для base_auth, либо интеграция с yandex Oauth2 для пользователя, отдает JWT токен, с которым пользователь может пользоваться другими сервисами
2. **UserService** – Управление пользователями. Создание, хранение, удаление.
3. **MenuService** – Управление меню.
4. **OrderService** – Управление заказами.
5. **PaymentService** – Обработка платежей (заглушка).
6. **Gateway** - начальная точка входа в приложение, роутинг, конвертация gRPC ответов от сервисов в Http ответы и обратно
7. **UserFrontend** - Web приложения для пользователя, заказы, просмотр меню, регистрация, просмотр статуса заказа, оплата()
8. **PersonalFrontend** - Web приложение для работников, оформление заказа, подтверждение, передача на кухню, отметка о готовности,

=== Структура базы данных

Каждый сервис имеет свою базу данных (PostgreSQL):

**AuthServiceDB:**

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    phone VARCHAR(20) UNIQUE,
    password_hash TEXT NOT NULL,
    full_name VARCHAR(100) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    last_login TIMESTAMPTZ
);
```

```sql
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
```

PostgreSQL + S3
В PostgreSQL рядом с информацией о меню/блюде хранится идентификатор на s3 по которому можно достать это фото и отдать на фронт
**MenuServiceDB:**

```sql
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,        -- Название категории: "Супы", "Десерты"
    description TEXT,                        -- Описание категории (опционально)
    display_order INT DEFAULT 0,             -- Порядок отображения в меню (0 = первый)
    is_active BOOLEAN DEFAULT TRUE           -- Активна ли категория (скрыть/показать)
);

CREATE INDEX idx_categories_active ON categories(is_active);
CREATE TABLE dishes (
    id SERIAL PRIMARY KEY,
    name VARCHAR(200) NOT NULL,              -- Название блюда: "Стейк Рибай"
    description TEXT,                        -- Описание состава и особенностей
    price NUMERIC(10, 2) NOT NULL
        CHECK (price > 0),                   -- Цена (например: 1500.50)
    category_id INT NOT NULL
        REFERENCES categories(id)
        ON DELETE RESTRICT,                  -- Запрет удаления категории с блюдами
    cooking_time_min INT
        CHECK (cooking_time_min > 0),        -- Время приготовления в минутах (опционально)
    image_url VARCHAR(500),                  -- Ссылка на фото в S3: "https://bucket.s3.amazonaws.com/dishes/123.jpg"
    is_available BOOLEAN DEFAULT TRUE,       -- Доступно ли для заказа
    calories INT,                            -- Ккал (опционально)
);

CREATE INDEX idx_dishes_category ON dishes(category_id);
CREATE INDEX idx_dishes_availability ON dishes(is_available);

CREATE TABLE promotions (
    id SERIAL PRIMARY KEY,
    dish_id INT NOT NULL REFERENCES dishes(id),
    discount_percent INT CHECK (discount_percent BETWEEN 1 AND 9900), -- Скидка 10%
    start_date DATE NOT NULL,
    end_date DATE NOT NULL,
    CONSTRAINT valid_dates CHECK (end_date > start_date)
);
CREATE TABLE menus (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,        -- "Основное меню", "Летнее спецпредложение"
    description TEXT,
    valid_from DATE NOT NULL,                 -- Дата начала действия меню
    valid_to DATE NOT NULL,                   -- Дата окончания
    is_active BOOLEAN DEFAULT TRUE,           -- Активно ли для отображения
    created_at TIMESTAMPTZ DEFAULT NOW(),
    CONSTRAINT valid_dates CHECK (valid_to >= valid_from)
);
-- Связь меню с категориями
CREATE TABLE menu_categories (
    menu_id INT NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    category_id INT NOT NULL REFERENCES categories(id),
    display_order INT DEFAULT 0,              -- Порядок категорий внутри меню
    PRIMARY KEY (menu_id, category_id)
);

-- Связь меню с отдельными блюдами (если нужно включать вне категорий)
CREATE TABLE menu_dishes (
    menu_id INT NOT NULL REFERENCES menus(id) ON DELETE CASCADE,
    dish_id INT NOT NULL REFERENCES dishes(id),
    PRIMARY KEY (menu_id, dish_id)
);
```

Эти сервисы будут на Clichouse
**OrderServiceDB:**

```sql
CREATE TABLE orders (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id INT NOT NULL,                    -- ID клиента (из UserService)
    staff_id INT,                            -- ID сотрудника (официанта, из StaffService)
    type VARCHAR(20) NOT NULL
        CHECK (type IN ('dine-in', 'delivery', 'takeaway')),
    status VARCHAR(50) NOT NULL
        CHECK (status IN ('created', 'confirmed', 'cooking', 'ready', 'delivered', 'canceled')),
    delivery_address JSONB,                  -- Для типа 'delivery' (город, улица, квартира)
    total_amount NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    finished_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_user ON orders(user_id);

CREATE TABLE order_items (
    order_id UUID NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    dish_id INT NOT NULL,                    -- ID блюда (из MenuService)
    quantity INT NOT NULL CHECK (quantity > 0),
    price NUMERIC(10, 2) NOT NULL,           -- Цена на момент заказа (фиксируется)
    special_requests TEXT,                   -- "Без лука, добавить соус"
    cooking_status VARCHAR(50)               -- Статус из KitchenService: 'pending', 'cooking', 'ready'
);

CREATE INDEX idx_order_items_dish ON order_items(dish_id);

CREATE TABLE order_status_history (
    order_id UUID NOT NULL REFERENCES orders(id),
    status VARCHAR(50) NOT NULL,
    changed_by INT,                          -- ID сотрудника или системы (NULL = автоматически)
    changed_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_status_history_order ON order_status_history(order_id);

CREATE TYPE payment_method AS ENUM ('cash', 'card', 'online');

CREATE TABLE payments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id UUID NOT NULL UNIQUE REFERENCES orders(id),
    amount NUMERIC(10, 2) NOT NULL CHECK (amount > 0),
    method payment_method NOT NULL,
    status VARCHAR(20) NOT NULL
        CHECK (status IN ('pending', 'completed', 'failed', 'refunded')),
    transaction_id VARCHAR(255),             -- ID транзакции в payment_service
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_payments_order ON payments(order_id);
```

**PaymentServiceDB:**

```sql
CREATE TABLE payments (
    id SERIAL PRIMARY KEY,
    order_id INT,
    amount DECIMAL(10,2),
    status VARCHAR(50),
    payment_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

== Implementation

### **Этапы разработки (2 месяца)**

#### **Неделя 1-2: Аутентификация и управление пользователями**

- Разработать **AuthService** (JWT, регистрация, вход).
- Реализовать **UserService** (CRUD пользователей).
- Подключить API Gateway для проверки токенов.
- Начать разработку UI (страницы входа/регистрации).

#### **Неделя 3-4: Управление меню**

- Разработать **MenuService** (CRUD блюд, категории).
- Подключить к фронтенду (страницы меню, редактирование).

#### **Неделя 5-6: Заказы и обработка событий**

- Разработать **OrderService** (создание/обновление заказов).
- Интеграция с Kafka для событий.
- Добавить интерфейсы заказов на фронтенде (клиент, официант, повар).

#### **Неделя 7: Платежи (заглушка) и финальная сборка**

- Разработать **PaymentService** (заглушка).
- Подключить к **OrderService** (изменение статуса заказа после оплаты).
- Завершить разработку UI (страница оплаты, статусы заказов).

#### **Неделя 8: Тестирование, отладка, развертывание**

- Интеграционное тестирование API и UI.
- Настроить мониторинг (Prometheus, Grafana, Loki).
- Настройка деплоя Terraform, Ansible конфигурации, настройка всех сервисов и бд.

== Gathering Results

### **Критерии успешности**

#### **1. Функциональность**

- Все основные сервисы работают корректно.
- Возможность создавать заказы, редактировать меню и управлять пользователями.
- Платёжная заглушка успешно меняет статус заказа.

#### **2. Производительность**

- Система выдерживает **500-1000 одновременных пользователей**.
- Среднее время отклика API **≤ 200 мс**.

#### **3. UI/UX**

- Интерфейс удобен для клиентов, официантов, поваров и администраторов.
- Минимальное количество кликов для оформления заказа.
- Поддержка мобильных устройств.

#### **4. Масштабируемость**

- Возможность интеграции с реальной платёжной системой.
- Готовность к добавлению доставки и других фич.

---
