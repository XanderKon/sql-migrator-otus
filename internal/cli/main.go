package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/XanderKon/sql-migrator-otus/internal/cli/command"
	"github.com/XanderKon/sql-migrator-otus/internal/cli/config"
	"github.com/XanderKon/sql-migrator-otus/internal/logger"
	"github.com/XanderKon/sql-migrator-otus/pkg/core"
)

var configFile string

func init() {
	if flag.Lookup("config") == nil {
		flag.StringVar(&configFile, "config", "./configs/config.yml", "Path to configuration file")
	}
}

func initFlags() {
	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr,
			`Usage: gomigrator [OPTIONS] COMMAND [arg...]
  
  You can override varuables from config file by ENV, just use something like "${DB_DSN}"

  OPTIONS:
    -config         Path to configuration file
		
  COMMAND:
    create [name]   Create migration with 'name'
    up              Migrate the DB to the most recent version available
    down            Roll back the version by 1
    redo            Re-run the latest migration
    status          Print all migrations status
    dbversion       Print migrations status (last applied migration)
    help            Print usage
    version         Application version

  Examples:
    gomigrator -config="../configs/config-test.yml" create "create_user_table"
    DB_DSN="postgresql://app:test@pgsql:5432/app" gomigrator up

Feel free to put PR here: https://github.com/XanderKon/sql-migrator-otus

Inspired by:
https://github.com/pressly/goose
https://github.com/golang-migrate/migrate
`)
	}
}

func Main() {
	initFlags()

	// init config
	cfg := config.NewConfig()

	// init logger
	logger := logger.New(cfg.Logger.Level, os.Stdout)

	// get args
	args := flag.Args()

	if len(args) == 0 {
		printUsage()
	}

	var cmd command.Command

	// init migrate api
	migrator, err := core.NewMigrator(cfg.Migrator.DSN, cfg.Migrator.TableName, cfg.Migrator.Dir)
	if err != nil {
		logger.Error("can't initialize migrator api! %s", err)
		return
	}

	// add logger
	migrator.Log = logger

	// close migrator (DB connection in simple case)
	defer migrator.Close()

	switch flag.Arg(0) {
	case "create":
		cmd = &command.Create{
			Cfg:    &cfg.Migrator,
			Logger: logger,
		}
	case "up":
		cmd = &command.Up{
			Migrator: migrator,
			Logger:   logger,
		}
	case "down":
		cmd = &command.Down{
			Migrator: migrator,
			Logger:   logger,
		}
	case "redo":
		cmd = &command.Redo{
			Migrator: migrator,
			Logger:   logger,
		}
	case "dbversion":
		cmd = &command.Dbversion{
			Migrator: migrator,
			Logger:   logger,
		}
	case "status":
		cmd = &command.Status{
			Migrator: migrator,
		}
	default:
		printUsage()
	}

	err = cmd.Run(args[1:])
	if errors.Is(err, core.ErrAlreadyUpToDate) || errors.Is(err, core.ErrNoAvailableMigrations) {
		logger.Info(err.Error())
	} else if err != nil {
		logger.Error("Error executing CLI: %s\n", err.Error())
		logger.Info("Try 'gomigrator help' for more information.")
	}
}

func printUsage() {
	flag.Usage()
	os.Exit(1)
}
