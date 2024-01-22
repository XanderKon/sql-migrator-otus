package parser

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testData struct {
	label string
	sql   string
	err   error
}

var sqlexample = `-- +gomigrator Up
CREATE TABLE IF NOT EXISTS test (
	id serial NOT NULL,
	test text
);
SELECT * FROM test;

-- +gomigrator Down
DROP TABLE test;
`

var errorexample = `-- CREATE TABLE IF NOT EXISTS test (
	id serial NOT NULL,
	test text
);
SELECT * FROM test;

-- +gomigrator Down
DROP TABLE test;
`

func TestParser(t *testing.T) {
	tests := []testData{
		{
			label: "success case",
			sql:   sqlexample,
			err:   nil,
		},
		{
			label: "error case",
			sql:   errorexample,
			err:   ErrIncorrectTemplate,
		},
	}

	for _, test := range tests {
		t.Run(test.label, func(t *testing.T) {
			migration, err := ParseMigration(strings.NewReader(test.sql))
			if test.err == nil {
				if err != nil {
					t.Fatalf("Unexpected error: %s", err.Error())
				}
			} else {
				assert.ErrorIs(t, err, test.err)
				return
			}

			// Not empty
			assert.NotEmpty(t, migration.UpStatements)
			assert.NotEmpty(t, migration.DownStatements)

			// Contains
			assert.Contains(t, migration.UpStatements, "CREATE TABLE IF NOT EXISTS test (\n\tid")
			assert.Contains(t, migration.DownStatements, "DROP TABLE test;")
		})
	}
}
