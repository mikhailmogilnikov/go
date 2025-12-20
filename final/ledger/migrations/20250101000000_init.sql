-- +goose Up
-- Таблица транзакций
CREATE TABLE IF NOT EXISTS transactions (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    amount NUMERIC(14,2) NOT NULL CHECK (amount > 0),
    category TEXT NOT NULL,
    description TEXT,
    date DATE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Индексы для транзакций
CREATE INDEX IF NOT EXISTS idx_transactions_user_id ON transactions(user_id);
CREATE INDEX IF NOT EXISTS idx_transactions_user_category ON transactions(user_id, category);
CREATE INDEX IF NOT EXISTS idx_transactions_user_date ON transactions(user_id, date);

-- Таблица бюджетов
CREATE TABLE IF NOT EXISTS budgets (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    category TEXT NOT NULL,
    limit_amount NUMERIC(14,2) NOT NULL CHECK (limit_amount > 0),
    period TEXT NOT NULL DEFAULT 'monthly',
    UNIQUE(user_id, category)
);

-- Индекс для бюджетов
CREATE INDEX IF NOT EXISTS idx_budgets_user_id ON budgets(user_id);

-- +goose Down
DROP INDEX IF EXISTS idx_budgets_user_id;
DROP TABLE IF EXISTS budgets;
DROP INDEX IF EXISTS idx_transactions_user_date;
DROP INDEX IF EXISTS idx_transactions_user_category;
DROP INDEX IF EXISTS idx_transactions_user_id;
DROP TABLE IF EXISTS transactions;



