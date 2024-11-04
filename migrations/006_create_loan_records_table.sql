-- migrations/006_create_loan_records_table.sql

CREATE TABLE IF NOT EXISTS loan_records (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL REFERENCES books(id),
    user_id UUID NOT NULL REFERENCES users(id),
    loan_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    due_date TIMESTAMPTZ NOT NULL,
    returned BOOLEAN DEFAULT FALSE,
    return_date TIMESTAMPTZ,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);