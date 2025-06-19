-- Создание таблицы сотрудников
CREATE TABLE IF NOT EXISTS staff (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    work_email VARCHAR(255) UNIQUE NOT NULL,
    work_phone VARCHAR(20) UNIQUE NOT NULL,
    full_name VARCHAR(255) NOT NULL,
    position VARCHAR(255) NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    hire_date TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- Создание индексов для таблицы сотрудников
CREATE INDEX IF NOT EXISTS idx_staff_work_email ON staff(work_email);
CREATE INDEX IF NOT EXISTS idx_staff_work_phone ON staff(work_phone);
CREATE INDEX IF NOT EXISTS idx_staff_is_active ON staff(is_active);
CREATE INDEX IF NOT EXISTS idx_staff_hire_date ON staff(hire_date);
CREATE INDEX IF NOT EXISTS idx_staff_position ON staff(position); 