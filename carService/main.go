package main

import (
	"amaksimov/carService/db"
	"amaksimov/carService/server"

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

	if err := server.CreateCarServer(); err != nil {
		log.Fatalf("error trying create car server: %v", err)
	}
}
