-- +gomigrator Up
SELECT * FROM test_table tt INNER JOIN some_table st ON tt.id = st.foreign_k WHERE tt.status = 1 ORDER BY tt.field LIMIT 15;

-- +gomigrator Down
TRUNCATE grass;
            