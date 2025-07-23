-- CREATE TYPE headlight_types AS ENUM('hid','led','tungsten','unspecified');
-- CREATE TYPE vehicle_sources AS ENUM('judicial','commission','overseas','unspecified');
-- CREATE TYPE wheel_sides AS ENUM('left','right','unspecified');
-- CREATE TYPE special_incident_types AS ENUM('casualty','suicide','watered','unspecified');

-- CREATE TABLE IF NOT EXISTS listed_vehicle (
--     listed_vehicle_id       SERIAL                  PRIMARY KEY,
--     manufactured_at         timestamptz             NOT NULL,
--     licensed_at             timestamptz,    
--     model_id                int                     NOT NULL REFERENCES vehicle_models(vehicle_model_id),
--     door_count              int                     NOT NULL,
--     key_count               int                     NOT NULL DEFAULT 0,
--     color                   text                    NOT NULL,
--     fuel_type               fuel_types              NOT NULL,
--     body_serial             text                    NOT NULL DEFAULT '',   
--     transmission_type       transmission_types      NOT NULL DEFAULT 'unspecified', 
--     headlight_type          headlight_types         NOT NULL DEFAULT 'unspecified',    
--     wheel_side              wheel_sides             NOT NULL DEFAULT 'unspecified',
--     engine_serial           text                    NOT NULL DEFAULT '',   

--     body_modified           boolean                 NOT NULL DEFAULT false,
--     vehicle_source          vehicle_sources         NOT NULL DEFAULT 'unspecified',    

--     special_incident        special_incident_types  NOT NULL DEFAULT 'unspecified',
--     note                    text                    NOT NULL DEFAULT '',

--     created_at              timestamptz             NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     updated_at              timestamptz             NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     deleted_at              timestamptz             
-- );

