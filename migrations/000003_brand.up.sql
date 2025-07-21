CREATE TABLE IF NOT EXISTS brands (
    brand_id serial PRIMARY KEY,
    brand_name text NOT NULL check( brand_name <> ''),
    logo_image text ,
    created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamptz DEFAULT CURRENT_TIMESTAMP
);