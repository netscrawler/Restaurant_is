-- Заполнение ролей для доставки
INSERT INTO roles (role_name) VALUES
('CUSTOMER'),   -- Клиент
('COURIER'),    -- Курьер
('MANAGER'),    -- Менеджер
('ADMIN'),      -- Администратор системы
('CHEF');      -- Повар

-- Заполнение дефолтных значений
INSERT INTO account_types (type_name) VALUES
('USER'),   -- Обычный пользователь
('STAFF');  -- Сотрудник службы доставки
