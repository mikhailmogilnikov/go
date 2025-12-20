-- Создаём базы данных для сервисов
CREATE DATABASE auth;
CREATE DATABASE ledger;

-- Подключаемся к auth и создаём таблицы
\c auth;
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Подключаемся к ledger и создаём таблицы
\c ledger;
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount NUMERIC(14,2) NOT NULL CHECK (amount > 0),
    category TEXT NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_category ON transactions(user_id, category);
CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, date);

CREATE TABLE IF NOT EXISTS budgets (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    category TEXT NOT NULL,
    limit_amount NUMERIC(14,2) NOT NULL CHECK (limit_amount > 0),
    period TEXT NOT NULL DEFAULT 'monthly',
    UNIQUE(user_id, category)
);
CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets(user_id);



