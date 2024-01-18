package main

import (
	"flag"

	"github.com/XanderKon/sql-migrator-otus/internal/cli"
)

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	cli.Main()
}
