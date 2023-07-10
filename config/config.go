package config

import (
	"log"
	"os"

	"gopkg.in/go-ini/ini.v1"
)

type ConfigList struct {
	ApiKey string
	ApiSercet string
	Logfile string
}

var Config ConfigList

func init() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		log.Printf("Failed to read file: %v", err)
		os.Exit(1)
	}

	Config = ConfigList{
		ApiKey: cfg.Section("bitflyer").Key("api_key").String(),
		ApiSercet: cfg.Section("bitflyer").Key("api_secret").String(),
		Logfile: cfg.Section("gotrading").Key("log_file").String(),
	}
}