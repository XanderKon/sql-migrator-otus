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
}

func Main() {
	flag.Parse()
	initFlags()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()

	logg := logger.New(config.Logger.Level, os.Stdout)

	logg.Info("Start migrator app")

	fmt.Println("Configuration: ", config.Migrator.DSN, config.Migrator.Dir, config.Migrator.Type, config.Logger.Level)

}
