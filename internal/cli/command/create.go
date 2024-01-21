package command

import (
	"errors"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/XanderKon/sql-migrator-otus/internal/cli/config"
	"github.com/XanderKon/sql-migrator-otus/internal/logger"
)

var ErrMissingName = errors.New("no migration name was set")

type tmplVars struct {
	Name string
}

type Create struct {
	Cfg    *config.MigratorConf
	Logger *logger.Logger
}

func (c *Create) Run(args []string) error {
	if len(args) == 0 {
		return ErrMissingName
	}

	return c.create(args[0])
}

func (c *Create) create(name string) error {
	// define new version of migration file.
	version := time.Now().UTC().UnixMilli()

	// define full filename
	fullname := fmt.Sprintf("%v_%v", version, c.snakeCase(name))
	filename := fmt.Sprintf("%v.%v", fullname, c.Cfg.Type)

	// define template
	var tmpl *template.Template
	if c.Cfg.Type == "go" {
		tmpl = goMigrationTemplate
	} else {
		tmpl = sqlMigrationTemplate
	}

	// try to create path
	err := os.MkdirAll(c.Cfg.Dir, 0o755)
	if err != nil {
		return fmt.Errorf("failed to create migration folder: %w", err)
	}

	// target path
	path := filepath.Join(c.Cfg.Dir, filename)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return fmt.Errorf("failed to create migration file: %w", err)
	}

	// try to create file
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create migration file2: %w", err)
	}
	defer f.Close()

	// compile template
	vars := tmplVars{
		Name: fullname,
	}
	if err := tmpl.Execute(f, vars); err != nil {
		return fmt.Errorf("failed to execute tmpl: %w", err)
	}

	c.Logger.Info("Success create new migration %s", filename)
	return nil
}

func (c *Create) snakeCase(s string) string {
	var b strings.Builder

	diff := 'a' - 'A'
	l := len(s)
	for i, v := range s {
		// replace all dots and other "danger" symbols
		ss := string(v)
		if ss == "+" || ss == "-" || ss == "—" || ss == "." || ss == "/" || ss == "," || ss == "_" {
			b.WriteRune('_')
			continue
		}

		if v >= 'a' {
			b.WriteRune(v)
			continue
		}
		if (i != 0 || i == l-1) && ((i > 0 && rune(s[i-1]) >= 'a') ||
			(i < l-1 && rune(s[i+1]) >= 'a')) {
			b.WriteRune('_')
		}
		b.WriteRune(v + diff)
	}
	return b.String()
}

var sqlMigrationTemplate = template.Must(template.New("gomigrator.sql-migration").Parse(`-- +gomigrator Up
SELECT 'up SQL query';

-- +gomigrator Down
SELECT 'down SQL query';
`))

var goMigrationTemplate = template.Must(template.New("gomigrator.go-migration").Parse(`package migrations

import (
	"context"
	"database/sql"
)

func Up_{{.Name}}(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is applied.
	return nil
}

func Down_{{.Name}}(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
`))
