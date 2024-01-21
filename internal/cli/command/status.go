package command

import (
	"errors"
	"os"

	"github.com/XanderKon/sql-migrator-otus/pkg/core"
	"github.com/jedib0t/go-pretty/v6/table"
)

var ErrGeneralError = errors.New("unable to show status table")

type Status struct {
	Migrator *core.Migrate
}

func (c *Status) Run(_ []string) error {
	migrations, err := c.Migrator.FullList()
	if err != nil {
		return ErrGeneralError
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Version", "Name", "Applied At"})

	for i, migr := range migrations {
		t.AppendRows([]table.Row{
			{i + 1, migr.Version, migr.Source, migr.AppliedAt.Format("2006-01-02 15:04:05")},
		})
	}

	t.AppendSeparator()
	t.AppendFooter(table.Row{"", "Total", len(migrations)})
	t.Render()

	return nil
}
