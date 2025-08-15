CREATE TABLE IF NOT EXISTS banner_stats (
    ts_minute TIMESTAMP NOT NULL,
    banner_id TEXT NOT NULL,
    count INT NOT NULL DEFAULT 0,
    PRIMARY KEY (ts_minute, banner_id)
);