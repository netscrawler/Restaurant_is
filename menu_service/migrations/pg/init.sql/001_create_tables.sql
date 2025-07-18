CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,        -- Название категории: "Супы", "Десерты"
    description TEXT,                        -- Описание категории (опционально)
    is_active BOOLEAN DEFAULT TRUE           -- Активна ли категория (скрыть/показать)
);

CREATE INDEX idx_categories_active ON categories(is_active);

CREATE TABLE dishes (
    id UUID PRIMARY KEY,
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

    created_at TIMESTAMP NOT NULL DEFAULT now(),  -- Дата создания
    updated_at TIMESTAMP NOT NULL DEFAULT now()   -- Дата последнего обновления
);


CREATE INDEX idx_dishes_category ON dishes(category_id);
CREATE INDEX idx_dishes_availability ON dishes(is_available);

INSERT INTO categories (name, description, is_active) VALUES
('Супы', 'Горячие первые блюда', TRUE),
('Салаты', 'Свежие и питательные салаты', TRUE),
('Горячие блюда', 'Основные блюда из мяса, рыбы и овощей', TRUE),
('Гарниры', 'Дополнения к основным блюдам', TRUE),
('Десерты', 'Сладкие блюда и выпечка', TRUE),
('Напитки', 'Безалкогольные напитки', TRUE),
('Завтраки', 'Блюда для утреннего меню', TRUE),
('Пасты', 'Итальянская паста и соусы', TRUE),
('Пицца', 'Разнообразные пиццы', TRUE),
('Соусы', 'Дополнительные соусы к блюдам', TRUE);

