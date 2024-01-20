package core

import (
	"context"
	"database/sql"
)

type Migration struct {
	// Migration version
	Version string

	// Type of migration (sql or go)
	Type string

	// Path to file
	Source string

	// Link to next Migration
	Next *Migration

	// Link to prev Migration
	Prev *Migration

	// Slice for statements to run up (used by SQL-migrations)
	UpSQL []string

	// Slice for statements to run down (used by SQL-migrations)
	DownSQL []string
}

func New() *Migration {
	return &Migration{}
}

// Up runs an up migration.
func (m *Migration) Up(db *sql.DB) error {
	// ctx := context.Background()
	return nil
}

// Down runs a down migration.
func (m *Migration) Down(db *sql.DB) error {
	// ctx := context.Background()
	return nil
}

// Internal logic of migration here
func (m *Migration) run(ctx context.Context, db *sql.DB, direction bool) error {
	return nil
}
