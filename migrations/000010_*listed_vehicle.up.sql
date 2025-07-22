CREATE TYPE headlight_types AS ENUM('hid','led','tungsten','others');
CREATE TYPE vehicle_sources AS ENUM('judicial','commission','overseas','unknown');
CREATE TYPE wheel_sides AS ENUM('left','right','others');
CREATE TYPE special_incident_types AS ENUM('casualty','suicide','watered');

CREATE TABLE IF NOT EXISTS listed_vehicle (
    listed_vehicle_id       SERIAL                  PRIMARY KEY,
    manufactured_at         timestamptz             NOT NULL,
    licensed_at             timestamptz,    
    brand_id                int                     REFERENCES brands(brand_id),
    model_id                int                     REFERENCES vehicle_models(vehicle_model_id),
    door_count              int                     NOT NULL,
    key_count               int                     DEFAULT 0,
    color                   text                    NOT NULL,
    fuel_type               fuel_types              NOT NULL,
    body_serial             text,   
    transmission_type       transmission_types, 
    headlight_type          headlight_types,    
    wheel_side              wheel_sides             NOT NULL DEFAULT 'left',

    engine_displacement     int,    
    valves_count            int,    
    engine_serial           text,   
    engine_type             engine_types,   
    has_turbo               boolean,    

    body_modified           boolean                 NOT NULL DEFAULT false,
    vehicle_source          vehicle_sources,    

    special_incident        special_incident_types,
    note                    text                    NOT NULL DEFAULT '',

    created_at              timestamptz             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at              timestamptz             NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at              timestamptz             
);

