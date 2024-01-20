package stub

import (
	"io"

	"github.com/XanderKon/sql-migrator-otus/internal/database"
)

type Stub struct {
	tablename string
	isLocked  bool
	version   int
	list      []int
}

// init itself
func init() {
	s := Stub{}
	database.Register("stub", &s)
}

func (p *Stub) Open(url string, tablename string) (database.Driver, error) {
	// create new DB instance
	instance := &Stub{
		tablename: tablename,
	}

	return instance, nil
}

func (p *Stub) Close() error {
	return nil
}

func (p *Stub) Lock() error {
	// TODO
	if p.isLocked {
		return database.ErrLocked
	}
	return nil
}

func (p *Stub) Unlock() error {
	// TODO
	return nil
}

func (p *Stub) Run(migration io.Reader) error {
	// TODO
	return nil
}

func (p *Stub) SetVersion(version int) error {
	p.version = version

	return nil
}

func (p *Stub) DeleteVersion(version int) error {
	return nil
}

// Version returns the currently active version.
// When no migration has been applied, it must return version -1.
func (p *Stub) Version() (int, error) {
	return p.version, nil
}

// List returns the slice of all apllied versions of migraions.
// When no migration has been applied, it must return empty slice.
func (p *Stub) List() ([]int, error) {
	return p.list, nil
}

// Create migrations table
func (p *Stub) PrepareTable() error {
	return nil
}
