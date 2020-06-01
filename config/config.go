package config

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Production       bool
	ConnectionString string
	DatabaseName     string
	Port             string
}

func GetConfig() *Configuration {
	file, err := os.Open(getConfigFile())
	if err != nil {
		panic(err)
	}
	decoder := json.NewDecoder(file)
	var cfg Configuration
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	setPort(&cfg)
	return &cfg
}

func getConfigFile() string {
	envVariable := os.Getenv("WEBAPI_ENVIRONMENT")
	if envVariable == "PROD" {
		return "config.production.json"
	} else {
		return "config.development.json"
	}
}

func setPort(cfg *Configuration) {
	if cfg.Production {
		cfg.Port = ":" + os.Getenv("HTTP_PLATFORM_PORT")
	} else {
		cfg.Port = ":8080"
	}
}
