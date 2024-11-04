-- migrations/005_create_loan_requests_table.sql

CREATE TABLE IF NOT EXISTS loan_requests (
    id SERIAL PRIMARY KEY,
    book_id INT NOT NULL REFERENCES books(id),
    user_id UUID NOT NULL REFERENCES users(id),
    request_time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    status VARCHAR(20) NOT NULL CHECK (status IN ('PENDING', 'APPROVED', 'REJECTED')),
    reject_reason TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);