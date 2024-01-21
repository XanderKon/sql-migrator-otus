package command

import (
	"github.com/XanderKon/sql-migrator-otus/internal/logger"
	"github.com/XanderKon/sql-migrator-otus/pkg/core"
)

type Dbversion struct {
	Migrator *core.Migrate
	Logger   *logger.Logger
}

func (c *Dbversion) Run(_ []string) error {
	return c.Migrator.Dbversion()
}
