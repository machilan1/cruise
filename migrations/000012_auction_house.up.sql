CREATE TABLE IF NOT EXISTS auction_houses(
    auction_house_id        SERIAL             PRIMARY KEY,
    auction_house_name      text               NOT NULL CHECK(auction_house_name <> ''),
    address_detail          text               NOT NULL CHECK(address_detail <> ''),
    town_id                 int                REFERENCES city_towns(town_id),
    created_at              timestamptz        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              timestamptz        NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at              timestamptz
);