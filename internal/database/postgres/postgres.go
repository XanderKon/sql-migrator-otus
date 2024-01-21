package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/XanderKon/sql-migrator-otus/internal/database"
	// Dynamic build.
	_ "github.com/lib/pq"
)

var ErrConnClose = fmt.Errorf("can't close connection")

type Postgres struct {
	db        *sql.DB
	tablename string
	ctx       context.Context
	isLocked  bool
}

// init itself.
func init() {
	psql := Postgres{}
	database.Register("postgres", &psql)
	database.Register("postgresql", &psql)
}

func (p *Postgres) Open(url string, tablename string) (database.Driver, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	err = db.PingContext(ctx)
	if err != nil {
		return nil, err
	}

	// create new DB instance
	instance := &Postgres{
		db:        db,
		tablename: tablename,
		ctx:       ctx,
	}

	return instance, nil
}

func (p *Postgres) Close() error {
	if err := p.db.Close(); err != nil {
		return fmt.Errorf("conn close error: %w", err)
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

// just run migration statement in transactions mode.
func (p *Postgres) Run(migration io.Reader) error {
	migr, err := io.ReadAll(migration)
	if err != nil {
		return err
	}

	query := string(migr)
	if strings.TrimSpace(query) == "" {
		return nil
	}

	tx, err := p.db.BeginTx(p.ctx, &sql.TxOptions{})
	if err != nil {
		return err
	}

	if _, err := tx.Exec(query); err != nil {
		if errRollback := tx.Rollback(); errRollback != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (p *Postgres) SetVersion(version int64) error {
	const query = `
		INSERT INTO %s (version, applied_at)
		VALUES (%d, $1)
	`
	_, err := p.db.ExecContext(
		p.ctx,
		fmt.Sprintf(query, p.tablename, version),
		time.Now(),
	)

	return err
}

func (p *Postgres) DeleteVersion(version int64) error {
	const query = `DELETE FROM %s WHERE version = %d;`

	_, err := p.db.ExecContext(
		p.ctx,
		fmt.Sprintf(query, p.tablename, version),
	)

	return err
}

// Version returns the currently active version.
// When no migration has been applied, it must return version -1.
func (p *Postgres) Version() (int64, error) {
	const query = `SELECT version FROM %s ORDER BY version DESC LIMIT 1;`

	row := p.db.QueryRowContext(
		p.ctx,
		fmt.Sprintf(query, p.tablename),
	)

	var version int64
	err := row.Scan(
		&version,
	)

	// If not migrations applied yet
	if errors.Is(err, sql.ErrNoRows) {
		return -1, nil
	}

	// Some sql errors
	if err != nil {
		return -1, err
	}

	return version, nil
}

// List returns the slice of all apllied versions of migraions.
// When no migration has been applied, it must return empty slice.
func (p *Postgres) List() ([]int64, error) {
	const query = `SELECT version FROM %s ORDER BY version;`

	rows, err := p.db.QueryContext(p.ctx, fmt.Sprintf(query, p.tablename))
	if err != nil {
		return []int64{}, err
	}
	defer rows.Close()

	versions := make([]int64, 0)

	for rows.Next() {
		var version int64
		err := rows.Scan(
			&version,
		)
		if err != nil {
			return nil, err
		}

		versions = append(versions, version)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return versions, nil
}

// Create migrations table.
func (p *Postgres) PrepareTable() error {
	const query = `
		CREATE TABLE IF NOT EXISTS %s (
			id serial NOT NULL,
			version bigint NOT NULL,
			applied_at timestamp NOT NULL,
			PRIMARY KEY(id),
			UNIQUE(version)
	);`
	_, err := p.db.ExecContext(
		p.ctx,
		fmt.Sprintf(query, p.tablename),
	)
	if err != nil {
		return err
	}

	return nil
}
