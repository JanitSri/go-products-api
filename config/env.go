package config

import (
	"log"

	"github.com/spf13/viper"
)

func GetEnvVaraiable(key string) string {

	viper.SetConfigFile("config/.env")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading in config file %s", err)
	}

	value, ok := viper.Get(key).(string)
	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}
