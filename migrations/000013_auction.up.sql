CREATE TABLE IF NOT EXISTS auctions(
    auction_id          bigserial            PRIMARY KEY,
    auction_date        timestamptz          NOT NULL,
    note                text,
    created_at          timestamptz          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at          timestamptz          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at          timestamptz          
);

