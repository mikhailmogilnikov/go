-- +goose Up
-- Create budgets table
CREATE TABLE IF NOT EXISTS budgets (
	id SERIAL PRIMARY KEY,
	category TEXT UNIQUE NOT NULL,
	limit_amount NUMERIC(14,2) NOT NULL CHECK (limit_amount > 0)
);

-- Create expenses table
CREATE TABLE IF NOT EXISTS expenses (
	id SERIAL PRIMARY KEY,
	amount NUMERIC(14,2) NOT NULL CHECK (amount <> 0),
	category TEXT NOT NULL,
	description TEXT,
	date DATE NOT NULL
);

-- Optional helpful index
CREATE INDEX IF NOT EXISTS idx_expenses_category_date ON expenses (category, date);

-- +goose Down
-- Drop optional index first (if exists)
DROP INDEX IF EXISTS idx_expenses_category_date;

-- Drop tables in reverse dependency order
DROP TABLE IF EXISTS expenses;
DROP TABLE IF EXISTS budgets;


