CREATE TABLE IF NOT EXISTS TABLE auction_house(
    auction_house_id SERIAL PRIMARY KEY,
    auction_house_name text NOT NULL CHECK(auction_house_name <> ''),
    address text NOT NULL CHECK(address <> ''),
    
)