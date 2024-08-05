CREATE TABLE IF NOT EXISTS strategy (
    id SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(20) NOT NULL,
    description TEXT NOT NULL,
    type VARCHAR(20) NOT NULL,
    mode VARCHAR(10) NOT NULL,
    mode_details TEXT NOT NULL,
    condition TEXT,
    query TEXT,
    status VARCHAR(10) NOT NULL,
    created_by VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_updated_by VARCHAR(30) NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    UNIQUE(name)
);

CREATE TABLE IF NOT EXISTS strategy_scrip (
    strategy_id INT NOT NULL,
    exchange VARCHAR(20) NOT NULL,
    token_type VARCHAR(20) NOT NULL,
    token VARCHAR(20) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    timestamp_epoch_millis bigint NOT NULL,
    created_by VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_updated_by VARCHAR(30) NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(strategy_id, exchange, token_type, token, transaction_type, timestmap_epoch_millis)
);

CREATE TABLE IF NOT EXISTS strategy_scrip_history (
    id SERIAL PRIMARY KEY NOT NULL,
    strategy_id INT NOT NULL,
    exchange VARCHAR(20) NOT NULL,
    token_type VARCHAR(20) NOT NULL,
    token VARCHAR(20) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    timestamp_epoch_millis bigint NOT NULL,
    created_by VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS strategy_trade (
    id SERIAL PRIMARY KEY NOT NULL,
    strategy_id INT NOT NULL,
    exchange VARCHAR(20) NOT NULL,
    token_type VARCHAR(20) NOT NULL,
    token VARCHAR(20) NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    quantity INT NOT NULL,
    meta VARCHAR(20) NOT NULL,
    created_by VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_updated_by VARCHAR(30) NOT NULL,
    last_updated_at TIMESTAMP NOT NULL
);

CREATE TABLE IF NOT EXISTS strategy_screener (
    strategy_id INT NOT NULL,
    screener_id INT NOT NULL,
    transaction_type VARCHAR(20) NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_by VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_updated_by VARCHAR(30) NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(strategy_id, screener_id)
);

CREATE TABLE IF NOT EXISTS strategy_backtest (
    strategy_id INT NOT NULL,
    backtest_id INT NOT NULL,
    status VARCHAR(20) NOT NULL,
    created_by VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_updated_by VARCHAR(30) NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(strategy_id, backtest_id)
);

CREATE TABLE IF NOT EXISTS scrip_delivery_traded (
    exchange VARCHAR(20) NOT NULL,
    token_type VARCHAR(20) NOT NULL,
    token VARCHAR(20) NOT NULL,
    trade_ts TIMESTAMP NOT NULL,
    traded_qty INT NOT NULL,
    delivered_qty INT NOT NULL,
    created_by VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    primary key(exchange, token_type, token, trade_ts)
);

CREATE TABLE IF NOT EXISTS scrip_analysis (
    exchange VARCHAR(20) NOT NULL,
    token_type VARCHAR(20) NOT NULL,
    token VARCHAR(20) NOT NULL,
    sentiment_score float8,
    created_by VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_updated_by VARCHAR(30) NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    primary key(exchange, token_type, token)
);

CREATE TABLE IF NOT EXISTS strategy_backtest_summary (
    strategy_id INT NOT NULL,
    backtest_id INT NOT NULL,
    exchange VARCHAR(20) NOT NULL,
    token_type VARCHAR(20) NOT NULL,
    token VARCHAR(20) NOT NULL,
    summary TEXT,
    status VARCHAR(20),
    created_by VARCHAR(30) NOT NULL,
    created_at TIMESTAMP NOT NULL,
    last_updated_by VARCHAR(30) NOT NULL,
    last_updated_at TIMESTAMP NOT NULL,
    PRIMARY KEY(strategy_id, backtest_id, exchange, token_type, token)
);

CREATE TABLE IF NOT EXISTS nifty_top_holdings (
    exchange VARCHAR(20) NOT NULL,
    token_type VARCHAR(20) NOT NULL,
    token VARCHAR(20) NOT NULL,
    weightage DOUBLE PRECISION NOT NULL,
    created_at  TIMESTAMP NOT NULL,
)

CREATE TABLE IF NOT EXISTS nse_historic_data (
    exchange VARCHAR(20) NOT NULL,
    token_type VARCHAR(20) NOT NULL,
    token VARCHAR(20) NOT NULL,
    industry VARCHAR(40) NOT NULL,
    weightage DOUBLE PRECISION NOT NULL,
    equity_capital DOUBLE PRECISION NOT NULL,
    free_float_market_capitalisation DOUBLE PRECISION NOT NULL,
    beta DOUBLE PRECISION NOT NULL,
    r2 DOUBLE PRECISION NOT NULL,
    volatility DOUBLE PRECISION NOT NULL,
    monthly_return_avg_impact_cost DOUBLE PRECISION NOT NULL,
    fetched_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL,
    PRIMARY KEY(exchange, token_type, token, created_at)
)
