CREATE TYPE drive_types AS ENUM('fwd','rwd','4wd','awd','unspecified');
CREATE TYPE fuel_types AS ENUM('diesel','gasoline','electric','gas','unspecified');
CREATE TYPE body_styles AS ENUM('sedan','wagon', 'hatchback','gt','sports','van','truck','suv','convertible','unspecified');
CREATE TYPE transmission_types AS ENUM('automatic','manual','unspecified');
CREATE TYPE engine_types AS ENUM('v','inline','boxer','rotary','unspecified');


CREATE TABLE IF NOT EXISTS series_variants(
    series_variant_id         serial              PRIMARY KEY,
    series_variant_name       text                NOT NULL CHECK(series_variant_name <> ''),
    version                 text                NOT NULL DEFAULT '',
    model_year              int                 NOT NULL,
    body_style              body_styles         NOT NULL DEFAULT 'unspecified',
    drive_type              drive_types         NOT NULL DEFAULT 'unspecified',
    fuel_type               fuel_types          NOT NULL DEFAULT 'unspecified',
    engine_type             engine_types        NOT NULL DEFAULT 'unspecified',
    engine_displacement     int                 NOT NULL CHECK(engine_displacement > 0),
    valve_count             int                 NOT NULL,
    has_turbo               bool                NOT NULL,
    transmission_type       transmission_types  NOT NULL DEFAULT 'unspecified',
    horse_power             int                 NOT NULL CHECK(horse_power > 0),
    series_id               int                 NOT NULL REFERENCES brand_series(brand_series_id),
    created_at              timestamptz         NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              timestamptz         NOT NULL DEFAULT CURRENT_TIMESTAMP
);

