-- +goose Up
-- +goose StatementBegin
CREATE TABLE users (
    chat_id BIGINT PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    name VARCHAR(255)
);
CREATE TABLE records (
    id UUID PRIMARY KEY,
    user_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL,
    type CHAR(4) NOT NULL,
    usdt_amount DOUBLE PRECISION NOT NULL,
    fiat_amount DOUBLE PRECISION NOT NULL,
    fiat_currency CHAR(3) NOT NULL,
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(chat_id)
        ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE records;
DROP TABLE users;
-- +goose StatementEnd
