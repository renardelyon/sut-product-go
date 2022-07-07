package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() (config Config, err error) {
	viper.AddConfigPath("./config/env")
	viper.SetConfigName("dev")
	viper.SetConfigType("env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		log.Println("Error: ", err.Error())
		return
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		log.Println("Error: ", err.Error())
		return
	}

	return
}
