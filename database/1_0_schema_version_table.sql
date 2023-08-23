-- Add schema changes table
CREATE TABLE schema_updates(
    id SERIAL NOT NULL PRIMARY KEY,
    major_version INT NOT NULL,
    minor_version INT NOT NULL,
    created_at TIMESTAMP NOT NULL
);

INSERT INTO schema_changes(major_version, minor_version) VALUES (1, 0, NOW());