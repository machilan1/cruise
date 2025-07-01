CREATE TABLE IF NOT EXISTS employees
(
    employee_id serial PRIMARY KEY,
    first_name  text NOT NULL CHECK (first_name <> ''),
    last_name   text NOT NULL CHECK (last_name <> ''),
    national_id text NOT NULL CHECK (national_id <> ''),
    email       text NOT NULL CHECK (email <> ''),
    cellphone   text NOT NULL CHECK (cellphone ~ '^\d{10}$'),
    created_at  timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at  timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE user_types AS ENUM ('superAdmin', 'admin', 'staff', 'user');
CREATE TABLE IF NOT EXISTS users
(
    user_id      serial NOT NULL PRIMARY KEY,
    user_type    user_types NOT NULL DEFAULT 'user',
    username     citext NOT NULL UNIQUE,
    display_name text NOT NULL CHECK (display_name <> ''),
    email        text UNIQUE,
    employee_id  int REFERENCES employees (employee_id) ON DELETE SET NULL,
    created_at   timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,

    email_quota int NOT NULL DEFAULT 3,
    CONSTRAINT email_quota_check CHECK (email_quota >= 0)
);

CREATE TABLE IF NOT EXISTS user_passwords
(
    user_password_id serial NOT NULL PRIMARY KEY,
    user_id          int NOT NULL REFERENCES users (user_id) ON DELETE CASCADE,
    password         bytea NOT NULL,
    created_at       timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS sessions
(
    session_id bytea NOT NULL PRIMARY KEY,
    data       bytea NOT NULL,
    user_id    int REFERENCES users (user_id) ON DELETE CASCADE,
    created_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);