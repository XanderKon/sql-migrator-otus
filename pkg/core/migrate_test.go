package core

import (
	"testing"

	_ "github.com/XanderKon/sql-migrator-otus/internal/database/stub"
	"github.com/stretchr/testify/assert"
)

var testMigrator *Migrate

func init() {
	testMigrator, _ = NewMigrator("stub://stub:stub@localhost:111/gomigrator", "migrations", "./test/migrations")
}

func TestFindAvailableMigrations(t *testing.T) {
	migrations, err := testMigrator.findAvailableMigrations()

	// "finder" works without problems
	assert.NoError(t, err)

	// exact two file exist
	assert.Len(t, migrations, 2)

	// check parse Version
	m := migrations[0]
	assert.Equal(t, "20240120195817", m.Version)

	m = migrations[1]
	assert.Equal(t, "20240120196753", m.Version)
}
