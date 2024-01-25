package main

import (
	"os"

	"github.com/XanderKon/sql-migrator-otus/internal/cli"
)

func main() {
	args := os.Args
	if len(args) > 1 && args[1] == "version" {
		printVersion()
		return
	}

	cli.Main()
}
