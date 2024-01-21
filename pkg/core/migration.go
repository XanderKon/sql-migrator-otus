package core

import (
	"database/sql"
	"time"
)

type Migration struct {
	// Migration version
	Version int64

	// Type of migration (sql or go)
	Type string

	// Path to file
	Source string

	// Path to file
	AppliedAt time.Time

	// Link to next Migration
	Next *Migration

	// Link to prev Migration
	Prev *Migration

	// Statements to run up (used by SQL-migrations)
	UpSQL string

	// Statements to run down (used by SQL-migrations)
	DownSQL string
}

func New() *Migration {
	return &Migration{}
}

// Up runs an up migration.
func (m *Migration) Up(_ *sql.DB) error {
	// ctx := context.Background()
	return nil
}

// Down runs a down migration.
func (m *Migration) Down(_ *sql.DB) error {
	// ctx := context.Background()
	return nil
}

// Internal logic of migration here.
// func (m *Migration) run(_ context.Context, _ *sql.DB, _ bool) error {
// 	return nil
// }
