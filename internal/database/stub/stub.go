package stub

import (
	"io"

	"github.com/XanderKon/sql-migrator-otus/internal/database"
)

type Stub struct {
	url       string
	tablename string
	isLocked  bool
	version   int
	list      []int
}

// init itself.
func init() {
	s := Stub{}
	database.Register("stub", &s)
}

func (p *Stub) Open(url string, tablename string) (database.Driver, error) {
	// create new DB instance
	instance := &Stub{
		url:       url,
		tablename: tablename,
	}

	return instance, nil
}

func (p *Stub) Close() error {
	return nil
}

func (p *Stub) Lock() error {
	if p.isLocked {
		return database.ErrLocked
	}
	return nil
}

func (p *Stub) Unlock() error {
	return nil
}

func (p *Stub) Run(_ io.Reader) error {
	return nil
}

func (p *Stub) SetVersion(version int) error {
	p.version = version

	return nil
}

func (p *Stub) DeleteVersion(_ int) error {
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

// Create migrations table.
func (p *Stub) PrepareTable() error {
	return nil
}
