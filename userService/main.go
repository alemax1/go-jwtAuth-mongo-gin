package main

import (
	"amaksimov/userService/db"

	"amaksimov/userService/server"
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error trying read from config: %v", err)
	}

	if err := db.CreateConnection(); err != nil {
		log.Fatalf("error trying connect to database: %v", err)
	}

	if err := server.CreateUserServer(); err != nil {
		log.Fatalf("error trying create user server: %v", err)
	}
}
