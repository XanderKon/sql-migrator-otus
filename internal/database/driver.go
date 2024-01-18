package database

import (
	"fmt"
	"io"
	"strings"
	"sync"
)

var (
	ErrParseDSN      = fmt.Errorf("can't get driver from dsn string")
	ErrUnknownDriver = fmt.Errorf("unknown driver")
	ErrLocked        = fmt.Errorf("can't acquire lock")
	ErrNotLocked     = fmt.Errorf("can't unlock, as not currently locked")
)

var driversMu sync.RWMutex

// list of available drivers for application
// [JUST POSTGRES IN OUR PROJECT CASE]
var drivers = make(map[string]Driver)

type Driver interface {
	// Open returns a new driver instance configured with parameters
	// coming from the URL string. Migrate will call this function
	// only once per instance.
	Open(url string) (Driver, error)

	// Close closes the underlying database instance managed by the driver.
	// Migrate will call this function only once per instance.
	Close() error

	// Lock should acquire a database lock so that only one migration process
	// can run at a time. Migrate will call this function before Run is called.
	// If the implementation can't provide this functionality, return nil.
	// Return database.ErrLocked if database is already locked.
	Lock() error

	// Unlock should release the lock. Migrate will call this function after
	// all migrations have been run.
	Unlock() error

	// Run applies a migration to the database. migration is guaranteed to be not nil.
	Run(migration io.Reader) error

	// SetVersion saves version.
	// Migrate will call this function before and after each call to Run.
	SetVersion(version int) error

	// DeleteVersion removes version.
	// Migrate will call this function before and after each call to Run.
	DeleteVersion(version int) error

	// Version returns the currently active version.
	// When no migration has been applied, it must return version -1.
	Version() (version int, err error)

	// List returns the slice of all aplloed migrations.
	// When no migration has been applied, it must return empty slice.
	List() (versions []int, err error)
}

// Register globally registers a driver.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()

	if driver == nil {
		panic("Register driver is nil")
	}

	if _, dup := drivers[name]; dup {
		panic("Register called twice for driver " + name)
	}

	drivers[name] = driver
}

// Open returns a new driver instance.
func Open(url string) (Driver, error) {
	i := strings.Index(url, ":")

	if i < 0 {
		return nil, ErrParseDSN
	}

	scheme := url[0:i]

	driversMu.RLock()
	d, ok := drivers[scheme]
	driversMu.RUnlock()
	if !ok {
		return nil, ErrUnknownDriver
	}

	return d.Open(url)
}
