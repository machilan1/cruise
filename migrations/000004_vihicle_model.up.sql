CREATE TYPE drive_types AS ENUM('fwd','rwd','4wd','awd','others');
CREATE TYPE fuel_types AS ENUM('diesel','gasoline','electric','gas','others');
CREATE TYPE body_styles AS ENUM('sedan','wagon', 'hatchback','gt','sports','van','truck','suv','others');
CREATE TYPE transmission_types AS ENUM('automatic','manual');
CREATE IF NOT EXISTS TYPE engine_types AS ENUM('v','inline','boxer','rotary','others');

CREATE TABLE IF NOT EXISTS vehicle_models (
    vehicle_model_id serial PRIMARY KEY,
    model_series_name text NOT NULL check (model_series_name <>''),
    model_commercial_name text NOT NULL check (model_commercial_name <>''),
    model_year int NOT NULL,
    brand_id int REFERENCES brands(brand_id),
    nickname text,
    engine_displacement int,
    drive_type drive_types NOT NULL,
    fuel_type fuel_types NOT NULL,
    body_style body_styles NOT NULL,
    transmission_type transmission_types NOT NULL,
    engine_type engine_types NOT NULL,
    CONSTRAINT uq_brand_id_series_name_model_year_model_commercial_name UNIQUE(brand_id,model_series_name,model_year,model_commercial_name)
)

-- <brand> <model_series_name> <model_year> <model_commercial_name> <engine_displacement>
-- example:
--  BMW 3-series 2014 320GT 2000
--  BMW 3-series 2014 320DGT 2000
--  BMW 3-series 2014 335 3000
--  MercedesBenz C-class 2020 C63 6300
--  Mitsubishi Lancer 2023 IO 1800
--  Mitsubishi Lancer 2023 Evolution 2500
-- Totota Yaris 2010 GLX 1600 


