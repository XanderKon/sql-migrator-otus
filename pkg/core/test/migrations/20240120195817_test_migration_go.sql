-- +gomigrator Up
SELECT * FROM migrations;
INSERT INTO migrations (version, applied_at) VALUES (2323, '2024-01-20 22:16:57.737429');

-- +gomigrator Down
SELECT 'down SQL query';
CREATE TABLE IF NOT EXISTS test (
			id serial NOT NULL,
			version bigint NOT NULL,
			applied_at timestamp NOT NULL,
			PRIMARY KEY(id),
			UNIQUE(version);
            