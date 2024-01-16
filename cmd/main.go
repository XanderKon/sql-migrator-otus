package main

import (
	"flag"
	"fmt"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "./configs/config.yml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config := NewConfig()

	fmt.Println("Configuration: ", config.DSN, config.Dir, config.Type)
}
