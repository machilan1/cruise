CREATE TABLE IF NOT EXISTS TABLE auction_houses(
    auction_house_id SERIAL PRIMARY KEY,
    auction_house_name text NOT NULL CHECK(auction_house_name <> ''),
    address text NOT NULL CHECK(address <> ''),
    town_id int REFERENCES towns(town_id),
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at timestamptz
)