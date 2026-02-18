CREATE TABLE loans (
    id              SERIAL PRIMARY KEY,
    principal       NUMERIC(15,2) NOT NULL,
    interest_rate   NUMERIC(5,2) NOT NULL,
    term_weeks      INT NOT NULL,
    weekly_amount   NUMERIC(15,2) NOT NULL,
    start_date      DATE NOT NULL,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE installments (
    id              SERIAL PRIMARY KEY,
    loan_id         INT NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    week_number     INT NOT NULL,
    due_date        DATE NOT NULL,
    amount          NUMERIC(15,2) NOT NULL,
    paid            BOOLEAN DEFAULT FALSE,
    UNIQUE(loan_id, week_number)
);

CREATE TABLE payments (
    id              SERIAL PRIMARY KEY,
    loan_id         INT NOT NULL REFERENCES loans(id) ON DELETE CASCADE,
    amount          NUMERIC(15,2) NOT NULL,
    payment_date    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    idempotency_key VARCHAR(255) UNIQUE
);

CREATE TABLE payment_installments (
    payment_id      INT NOT NULL REFERENCES payments(id) ON DELETE CASCADE,
    installment_id  INT NOT NULL REFERENCES installments(id) ON DELETE CASCADE,
    PRIMARY KEY (payment_id, installment_id)
);

CREATE INDEX idx_installments_loan_due ON installments(loan_id, due_date);
CREATE INDEX idx_installments_loan_paid ON installments(loan_id, paid);
CREATE INDEX idx_payments_idempotency ON payments(idempotency_key);