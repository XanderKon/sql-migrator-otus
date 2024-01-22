-- +gomigrator Up
INSERT INTO test (test, column_int, column_datetime) VALUES ('some test text', 1, '2024-01-20 22:16:57.737429');
INSERT INTO test (test, column_int, column_datetime) VALUES ('some another test text', 2, '2024-01-21 22:16:57.737429');

-- +gomigrator Down
TRUNCATE test;
