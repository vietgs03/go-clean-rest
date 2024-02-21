CREATE TABLE market_history (
    id SERIAL PRIMARY KEY,
    symbol VARCHAR(50) NOT NULL,
    time TIMESTAMP NOT NULL,
    open DOUBLE PRECISION NOT NULL,
    high DOUBLE PRECISION NOT NULL,
    low DOUBLE PRECISION NOT NULL,
    close DOUBLE PRECISION NOT NULL
);
CREATE INDEX idx_symbol ON market_history (symbol);
CREATE INDEX idx_time ON market_history (time);
