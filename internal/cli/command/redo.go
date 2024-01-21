package command

import (
	"github.com/XanderKon/sql-migrator-otus/internal/logger"
	"github.com/XanderKon/sql-migrator-otus/pkg/core"
)

type Redo struct {
	Migrator *core.Migrate
	Logger   *logger.Logger
}

func (c *Redo) Run(_ []string) error {
	return c.Migrator.Redo()
}
