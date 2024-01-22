package test

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/XanderKon/sql-migrator-otus/internal/database"
	"github.com/XanderKon/sql-migrator-otus/pkg/core"
	"github.com/stretchr/testify/suite"
)

type MigratorSuire struct {
	suite.Suite
	migrator *core.Migrate

	dsn    string
	driver database.Driver
}

const DefaultTableName = "migrations"

func (s *MigratorSuire) SetupSuite() {
	dsn := os.Getenv("DSN")
	dir := os.Getenv("DIR")

	s.dsn = dsn

	// init migrator by data from env.
	migrator, err := core.NewMigrator(dsn, DefaultTableName, dir)
	s.Require().NoError(err)
	s.Require().NotNil(migrator)
	s.migrator = migrator

	// init additional driver connection for checking.
	driver, err := database.Open(s.dsn, DefaultTableName)
	s.Require().NoError(err)
	s.Require().NotNil(driver)
	s.driver = driver
}

// close connection after finishing suite.
func (s *MigratorSuire) TearDownSuite() {
	defer s.migrator.Close()
}

// clear everything after each test.
func (s *MigratorSuire) TearDownTest() {
	query := fmt.Sprintf(`DROP TABLE IF EXISTS test;TRUNCATE %s`, DefaultTableName)
	err := s.driver.Run(strings.NewReader(query))
	s.Require().NoError(err)
}

func (s *MigratorSuire) TestMigratorUp() {
	// ensure that migrator exist.
	s.NotNil(s.T(), s.migrator)

	// run all up.
	err := s.migrator.Up()
	s.NoError(err)
	s.checkAppliedListCount(3)

	// run all up again.
	err = s.migrator.Up()
	s.ErrorIs(err, core.ErrAlreadyUpToDate)
	s.checkAppliedListCount(3)
}

func (s *MigratorSuire) TestMigratorDown() {
	// ensure that migrator exist
	s.NotNil(s.T(), s.migrator)

	// run all up.
	err := s.migrator.Up()
	s.NoError(err)
	s.checkAppliedListCount(3)

	// run one Down.
	err = s.migrator.Down()
	s.NoError(err)
	s.checkAppliedListCount(2)

	// one more Down.
	err = s.migrator.Down()
	s.NoError(err)
	s.checkAppliedListCount(1)

	// one more Down.
	err = s.migrator.Down()
	s.NoError(err)
	s.checkAppliedListCount(0)

	// final Down.
	err = s.migrator.Down()
	s.ErrorIs(err, core.ErrAlreadyUpToDate)
	s.checkAppliedListCount(0)
}

func (s *MigratorSuire) TestMigratorRedo() {
	// ensure that migrator exist.
	s.NotNil(s.T(), s.migrator)

	// run all up.
	err := s.migrator.Up()
	s.NoError(err)
	s.checkAppliedListCount(3)

	// get last migration.
	list, err := s.migrator.FullList()
	s.NoError(err)
	lastMigration := list[len(list)-1]

	// sleep...zzzz.
	time.Sleep(1 * time.Second)

	// run Redo
	err = s.migrator.Redo()
	s.NoError(err)
	s.checkAppliedListCount(3)

	// get last migration after Redo.
	list, err = s.migrator.FullList()
	s.NoError(err)
	lastMigrationAfterRedo := list[len(list)-1]

	s.NotEqual(lastMigration.AppliedAt.Unix(), lastMigrationAfterRedo.AppliedAt.Unix())
	s.Greater(lastMigrationAfterRedo.AppliedAt.Unix(), lastMigration.AppliedAt.Unix())
}

func (s *MigratorSuire) TestMigratorDbversion() {
	// ensure that migrator exist.
	s.NotNil(s.T(), s.migrator)

	// check empty version
	version, err := s.migrator.Dbversion()
	s.ErrorIs(err, core.ErrNoCurrentVersion)
	s.Equal(version, int64(-1))

	// run all up.
	err = s.migrator.Up()
	s.NoError(err)
	s.checkAppliedListCount(3)

	// get last migration.
	list, err := s.migrator.FullList()
	s.NoError(err)
	lastMigration := list[len(list)-1]

	// check that are equals
	version, err = s.migrator.Dbversion()
	s.NoError(err)
	s.Equal(lastMigration.Version, version)
}

// check applied list.
func (s *MigratorSuire) checkAppliedListCount(expectedCount int) {
	list, err := s.migrator.FullList()
	s.NoError(err)
	s.Equal(len(list), expectedCount)
}

func TestMigrator(t *testing.T) {
	suite.Run(t, new(MigratorSuire))
}
