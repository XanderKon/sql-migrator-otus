package command

import (
	"github.com/XanderKon/sql-migrator-otus/internal/logger"
	"github.com/XanderKon/sql-migrator-otus/pkg/core"
)

type Up struct {
	Migrator *core.Migrate
	Logger   *logger.Logger
}

func (c *Up) Run(_ []string) error {
	return c.Migrator.Up()
}
