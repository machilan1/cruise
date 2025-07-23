CREATE TABLE IF NOT EXISTS brand_series (
    brand_series_id         serial      PRIMARY KEY,
    brand_series_name       text        NOT NULL CHECK(brand_series_name <> ''),
    brand_id                int         REFERENCES brands(brand_id),
    created_at              timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT uq_brand_series_brand_id_brand_series_name UNIQUE(brand_series_name , brand_id)
)