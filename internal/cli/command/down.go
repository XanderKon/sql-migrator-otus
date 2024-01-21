package command

import (
	"github.com/XanderKon/sql-migrator-otus/internal/logger"
	"github.com/XanderKon/sql-migrator-otus/pkg/core"
)

type Down struct {
	Migrator *core.Migrate
	Logger   *logger.Logger
}

func (c *Down) Run(_ []string) error {
	return c.Migrator.Down()
}
