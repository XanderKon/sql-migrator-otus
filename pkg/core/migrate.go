package core

import (
	"errors"
	"fmt"
	"io/fs"
	"net/http"
	"path"
	"slices"
	"sort"
	"strconv"
	"strings"

	"github.com/XanderKon/sql-migrator-otus/internal/database"
	"github.com/XanderKon/sql-migrator-otus/internal/logger"
	"github.com/XanderKon/sql-migrator-otus/internal/parser"
)

var (
	ErrNoCurrentVersion      = errors.New("no current version found")
	ErrNoAvailableMigrations = errors.New("no available migrations found")
	ErrAlreadyUpToDate       = errors.New("already up to date")
)

type Migrate struct {
	Log logger.Logger

	driver    database.Driver
	tablename string
	dir       string
}

// Migrations slice.
type Migrations []*Migration

func NewMigrator(dsn string, tableName string, dir string) (*Migrate, error) {
	// get driver
	driver, err := database.Open(dsn, tableName)
	if err != nil {
		return nil, fmt.Errorf("can't get driver: %w", err)
	}

	migr := &Migrate{
		driver:    driver,
		tablename: tableName,
		dir:       dir,
	}

	// create table if does not exist
	err = migr.prepareDatabase()

	if err != nil {
		return nil, fmt.Errorf("can't initialize table: %w", err)
	}

	return migr, nil
}

func (m *Migrate) Up() error {
	if err := m.lock(); err != nil {
		return err
	}

	migrations, err := m.migrationsForRun(true, 0)

	if errors.Is(err, ErrNoAvailableMigrations) || errors.Is(err, ErrAlreadyUpToDate) {
		m.Log.Info(err.Error())
		return nil
	} else if err != nil {
		return err
	}

	for _, migr := range migrations {
		if err := m.driver.Run(strings.NewReader(migr.UpSQL)); err != nil {
			return fmt.Errorf("can't execute migration with version %d: %w", migr.Version, err)
		}

		// set version here if success
		m.setVersion(migr.Version)
		m.Log.Info("Migration %d successfully applied!", migr.Version)
	}

	defer m.unlock()
	return nil
}

func (m *Migrate) Down() {
}

func (m *Migrate) Redo() {
}

// close migrator API
// just close DB connection in our case.
func (m *Migrate) Close() error {
	return m.driver.Close()
}

// prepare migrations slice for next Run
// up -- direction
// limit -- how many migrations should be executed (0 -- without limit).
func (m *Migrate) migrationsForRun(up bool, _ int) (Migrations, error) {
	// get available migrations
	availableMigrations, err := m.findAvailableMigrations()
	if err != nil {
		return make(Migrations, 0), err
	}

	// no available migrations, so skip all the next
	if len(availableMigrations) == 0 {
		return make(Migrations, 0), ErrNoAvailableMigrations
	}

	// get list of applied migrations
	appliedVersions, err := m.list()
	if err != nil {
		return make(Migrations, 0), err
	}

	// if we go down and don't have any applied migrations - do nothing
	if len(appliedVersions) == 0 && !up {
		return make(Migrations, 0), ErrAlreadyUpToDate
	}

	// if we go up and don't have any applied migrations - run all of then
	if len(appliedVersions) == 0 && up {
		return availableMigrations, nil
	}

	var migrationsForRun Migrations

	// calc the difference between them
	if !up {
		// sort desc
		sort.Slice(availableMigrations, func(i, j int) bool {
			return availableMigrations[i].Version > availableMigrations[j].Version
		})

		// filter them
		for i, migr := range availableMigrations {
			if availableMigrations[i].Version <= appliedVersions[0] {
				migrationsForRun = append(migrationsForRun, migr)
			}
		}
	} else {
		// filter them
		for i, migr := range availableMigrations {
			// if it already applied - skip
			if slices.Contains(appliedVersions, availableMigrations[i].Version) {
				continue
			}

			if availableMigrations[i].Version > appliedVersions[0] {
				migrationsForRun = append(migrationsForRun, migr)
			}
		}
	}

	if len(migrationsForRun) == 0 {
		return make(Migrations, 0), ErrAlreadyUpToDate
	}

	return migrationsForRun, nil
}

// func (m *Migrate) getMigrationByVersion(migraions Migrations, version int64) (*Migration, error) {
// 	for _, migr := range migraions {
// 		if migr.Version == version {
// 			return migr, nil
// 		}
// 	}

// 	return nil, fmt.Errorf("no migration find by version %d", version)
// }

func (m *Migrate) findAvailableMigrations() (Migrations, error) {
	migrations := make([]*Migration, 0)

	file, err := http.Dir(m.dir).Open(".")
	if err != nil {
		return nil, err
	}

	files, err := file.Readdir(0)
	if err != nil {
		return nil, err
	}

	for _, info := range files {
		name := info.Name()
		if strings.HasSuffix(name, ".sql") {
			migration, err := m.parseSQLMigration(info)
			if err != nil {
				return nil, err
			}

			migrations = append(migrations, migration)
		}
	}

	// insure then they are sorted by version correctly
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

// parse SQL migration file.
func (m *Migrate) parseSQLMigration(info fs.FileInfo) (*Migration, error) {
	path := path.Join("./", info.Name())

	file, err := http.Dir(m.dir).Open(path)
	if err != nil {
		return nil, fmt.Errorf("error while opening %s: %w", info.Name(), err)
	}
	defer func() { _ = file.Close() }()

	if err != nil {
		return nil, fmt.Errorf("error while opening %s: %w", info.Name(), err)
	}

	version := m.getVersionFromFileName(info.Name())

	migration := &Migration{
		Version: version,
	}

	parsed, err := parser.ParseMigration(file)
	if err != nil {
		return nil, fmt.Errorf("error while parsing file %s: %w", info.Name(), err)
	}

	// set statements
	migration.UpSQL = parsed.UpStatements
	migration.DownSQL = parsed.DownStatements

	return migration, nil
}

func (m *Migrate) getVersionFromFileName(filename string) int64 {
	version := strings.Split(filename, "_")[0]

	i, _ := strconv.ParseInt(version, 10, 64)

	return i
}

// create migrations table if it doesn't exist.
func (m *Migrate) prepareDatabase() error {
	if err := m.lock(); err != nil {
		return err
	}

	defer m.unlock()
	return m.driver.PrepareTable()
}

func (m *Migrate) setVersion(version int64) error {
	if err := m.lock(); err != nil {
		return err
	}

	err := m.driver.SetVersion(version)
	if err != nil {
		return fmt.Errorf("can't set new migraion version: %w", err)
	}

	defer m.unlock()
	return nil
}

func (m *Migrate) list() ([]int64, error) {
	if err := m.lock(); err != nil {
		return []int64{}, err
	}

	list, err := m.driver.List()
	if err != nil {
		return []int64{}, fmt.Errorf("can't get list of applied migraions: %w", err)
	}

	defer m.unlock()

	return list, nil
}

// Get current migration version from DB driver.
// func (m *Migrate) current() (int, error) {
// 	if err := m.lock(); err != nil {
// 		return -1, err
// 	}

// 	curVersion, err := m.driver.Version()
// 	if err != nil {
// 		return -1, fmt.Errorf("can't get current migration: %w", err)
// 	}

// 	defer m.unlock()

// 	return curVersion, nil
// }

func (m *Migrate) lock() error {
	return m.driver.Lock()
}

func (m *Migrate) unlock() {
	if err := m.driver.Unlock(); err != nil {
		m.Log.Error(fmt.Errorf("can't unlock from database driver: %w", err).Error())
	}
}
