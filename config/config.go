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
	envVariable := os.Getenv("APPSETTINGS_ENVIRONMENT")
	if envVariable == "PROD" {
		return "config.production.json"
	} else {
		return "config.development.json"
	}
}

func setPort(cfg *Configuration) {
	port := ":" + os.Getenv("HTTP_PLATFORM_PORT")
	if port == ":" {
		port = ":8080"
	}
	cfg.Port = port
}
