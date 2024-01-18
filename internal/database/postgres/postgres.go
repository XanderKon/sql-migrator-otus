package postgres

import (
	"database/sql"
	"fmt"
	"io"

	"github.com/XanderKon/sql-migrator-otus/internal/database"
)

var (
	ErrConnClose = fmt.Errorf("can't close connection")
)

type Postgres struct {
	conn     *sql.Conn
	db       *sql.DB
	isLocked bool
}

// init itself
func init() {
	psql := Postgres{}
	database.Register("postgres", &psql)
}

func (p *Postgres) Open(url string) (database.Driver, error) {
	// TODO
	return nil, nil
}

func (p *Postgres) Close() error {
	connErr := p.conn.Close()
	var dbErr error
	if p.db != nil {
		dbErr = p.db.Close()
	}

	if connErr != nil || dbErr != nil {
		return fmt.Errorf("conn: %v, db: %v", connErr, dbErr)
	}
	return nil
}

func (p *Postgres) Lock() error {
	// TODO
	if p.isLocked {
		return database.ErrLocked
	}
	return nil
}

func (p *Postgres) Unlock() error {
	// TODO
	return nil
}

func (p *Postgres) Run(migration io.Reader) error {
	// TODO
	return nil
}

func (p *Postgres) SetVersion(version int) error {
	// TODO
	return nil
}

func (p *Postgres) DeleteVersion(version int) error {
	// TODO
	return nil
}

func (p *Postgres) Version() (version int, err error) {
	// TODO
	return 0, nil
}

func (p *Postgres) List() (versions []int, err error) {
	// TODO
	return make([]int, 0), nil
}
