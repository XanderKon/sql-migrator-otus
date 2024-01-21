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

	// check parse Version && correct order of files
	m := migrations[0]
	assert.Equal(t, int64(20240120195817), m.Version)

	m = migrations[1]
	assert.Equal(t, int64(20240120196753), m.Version)
}

func TestGetVersionFromFileName(t *testing.T) {
	version := testMigrator.getVersionFromFileName("1234_test_migr.sql")
	assert.NotEmpty(t, version)
	assert.Equal(t, version, int64(1234))
}
