package cli

import (
	"flag"
	"fmt"
	"os"

	"github.com/XanderKon/sql-migrator-otus/internal/logger"
)

var configFile string

func init() {
	if flag.Lookup("config") == nil {
		flag.StringVar(&configFile, "config", "./configs/config.yml", "Path to configuration file")
	}
}

func initFlags() {
	configFile = flag.Lookup("config").Value.(flag.Getter).Get().(string)

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

	config := NewConfig()

	logg := logger.New(config.Logger.Level, os.Stdout)

	logg.Info("Start migrator app")

	flag.Usage()
}
