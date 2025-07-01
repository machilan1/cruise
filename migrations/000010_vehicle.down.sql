-- CREATE IF NOT EXISTS TYPE vehicle_types AS ENUM('sedan','wagon','gt','sports','van','truck','suv','others');
-- CREATE IF NOT EXISTS TYPE transmission_types AS ENUM('automatic','manual');
-- CREATE IF NOT EXISTS TYPE fuel_types AS ENUM('diesel','gasoline','electric','gas','others');
-- CREATE IF NOT EXISTS TYPE headlight_types AS ENUM('hid','led','tungsten','others');
-- CREATE IF NOT EXISTS TYPE vehicle_sources AS ENUM('judicial','commission','oversea','unknown');
-- CREATE IF NOT EXISTS TYPE wheel_sides AS ENUM('left','right','others');
-- CREATE IF NOT EXISTS TYPE engine_types AS ENUM('v','inline','boxer','rotary','others');
-- CREATE IF NOT EXISTS TYPE drive_styles AS ENUM('fwd','rwd','4wd','awd';);
-- CREATE IF NOT EXISTS TYPE special_incident_types AS ENUM('casualty','suicide','watered');

-- CREATE IF NOT EXISTS TABLE vehicle (
--     vehicle_id SERIAL PRIMARY KEY,
--     vehicle_type vehicle_types NOT NULL,
--     model_year int,
--     brand_id int REFERENCES brands(brand_id),
--     model_id int,
--     manufactured_date date with time zone DEFAULT CURRENT_TIMESTAMP,
--     num_of_doors int NOT NULL,
--     num_of_keys int DEFAULT 0,
--     color text NOT NULL,
--     fuel_type fuel_types NOT NULL,
--     vim text,
--     transmission_type transmission_types,
--     headlight_type headlight_types,
--     wheel_side wheel_sides NOT NULL,
--     drive_style drive_styles,

--     engine_volume int,
--     engine_valves_num int,
--     engine_serial_num text,
--     engine_type engine_types,
--     turbo_charged boolean,
    
--     licensed_at timestamptz,
--     modified boolean NOT NULL DEFAULT false,
--     source vehicle_sources,
    
--     special_incident special_incident_types,
--     note text,

--     created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
--     updated_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP
-- )

