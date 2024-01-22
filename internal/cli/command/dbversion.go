package command

import (
	"errors"

	"github.com/XanderKon/sql-migrator-otus/internal/logger"
	"github.com/XanderKon/sql-migrator-otus/pkg/core"
)

type Dbversion struct {
	Migrator *core.Migrate
	Logger   *logger.Logger
}

func (c *Dbversion) Run(_ []string) error {
	_, err := c.Migrator.Dbversion()

	if errors.Is(err, core.ErrNoCurrentVersion) {
		c.Logger.Info(err.Error())
		return nil
	}

	return err
}
