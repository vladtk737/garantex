-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS exchange (
                                     id SERIAL PRIMARY KEY,
                                     market VARCHAR(50) NOT NULL,
                                     ask_price NUMERIC(20, 8) NOT NULL,
                                     bid_price NUMERIC(20, 8) NOT NULL,
                                     timestamp TIMESTAMPTZ NOT NULL,
                                     UNIQUE(market, timestamp)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS exchange;

-- +goose StatementEnd