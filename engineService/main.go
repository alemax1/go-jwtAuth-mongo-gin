package main

import (
	"amaksimov/engineService/cache"
	"amaksimov/engineService/db"
	"amaksimov/engineService/server"
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

	if err := server.CreateEngineServer(); err != nil {
		log.Fatalf("error trying create user server: %v", err)
	}
}
