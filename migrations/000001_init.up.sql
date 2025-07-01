-- 為加強view管理，下列嘗試管理政策如下
-- vb:  後端會實際使用的view，改欄位要注意
-- vm:  跟vb view欄位有不同，不使用於後端
-- vr:  用於報表的view，不使用於後端
-- vrt: 單純用於協助組成報表的view
-- vrb: 用於報表的view，且會使用於後端
-- vjb:  json object view，且會使用於後端
-- vj: json object view，不使用於後端

CREATE TABLE files
(
    file_id           serial PRIMARY KEY,
    path              text NOT NULL UNIQUE CHECK (path <> ''),
    original_filename text NOT NULL CHECK (original_filename <> ''),
    mime_type         text NOT NULL CHECK (mime_type <> ''),
    size              int NOT NULL CHECK (size >= 0),
    created_at        timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at        timestamp with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_files_path ON files (path);
CREATE INDEX idx2_files_path ON files (file_id, path);
CREATE INDEX idx3_files_ext ON files (mime_type);

-- below setting is to auto reset email quota without using crop job.
CREATE TYPE system_parameter_keys AS ENUM ('lastEmailRefreshDate');
CREATE TABLE system_parameters
(
    parameter_key system_parameter_keys PRIMARY KEY,
    value        text
);

CREATE VIEW vb_last_email_refresh_date AS(
    SELECT CAST(value AS date) AS last_email_refresh_date
    FROM system_parameters
    WHERE parameter_key = 'lastEmailRefreshDate'
);

INSERT INTO system_parameters (parameter_key, value) VALUES ('lastEmailRefreshDate', '2025-02-20');