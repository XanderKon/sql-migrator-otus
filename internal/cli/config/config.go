package config

import (
	"flag"
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	Migrator MigratorConf `mapstructure:"migrator"`
	Logger   LoggerConf   `mapstructure:"logger"`
}

type MigratorConf struct {
	DSN       string `mapstructure:"dsn"`
	Dir       string `mapstructure:"dir"`
	Type      string `mapstructure:"type"`
	TableName string `mapstructure:"table_name"`
}

type LoggerConf struct {
	Level string `mapstructure:"level"`
}

func NewConfig() *Config {
	configFile := flag.Lookup("config").Value.(flag.Getter).Get().(string)

	v := viper.New()
	v.AutomaticEnv()
	v.SetConfigFile(configFile)

	if err := v.ReadInConfig(); err != nil {
		fmt.Printf("couldn't load config: %s", err)
		os.Exit(1)
	}

	// Support to set ENV variables in file
	for _, k := range v.AllKeys() {
		val := v.GetString(k)
		v.Set(k, os.ExpandEnv(val))
	}

	var config Config
	if err := v.Unmarshal(&config); err != nil {
		fmt.Printf("couldn't read config: %s", err)
	}

	return &config
}
