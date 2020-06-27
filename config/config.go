package config

import (
	"encoding/json"
	"os"
	"strconv"
)

type Configuration struct {
	Production       string
	ConnectionString string
	DatabaseName     string
	SigningKey       string
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
	isProd, err := strconv.ParseBool(cfg.Production)
	if err != nil {
		panic(err)
	}
	if isProd {
		cfg.Port = ":" + os.Getenv("HTTP_PLATFORM_PORT")
	} else {
		cfg.Port = ":8080"
	}
}
