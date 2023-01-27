package main

import (
	"amaksimov/carService/cache"
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

	if err := cache.CreateRedisConnection(); err != nil {
		log.Fatalf("error creating redis connection: %v", err)
	}
	log.Println("Redis connection successfully created")

	if err := server.CreateCarServer(); err != nil {
		log.Fatalf("error trying create car server: %v", err)
	}
}
