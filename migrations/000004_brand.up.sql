CREATE TABLE IF NOT EXISTS brands (
    brand_id    serial          PRIMARY KEY,
    brand_name  text            NOT NULL CHECK( brand_name <> ''),
    image_id    int             REFERENCES files(file_id),
    created_at  timestamptz     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamptz     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_brands_brand_name UNIQUE(brand_name)
);