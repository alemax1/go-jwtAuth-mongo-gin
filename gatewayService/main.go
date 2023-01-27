package main

import (
	"amaksimov/gatewayService/grpc"
	"amaksimov/gatewayService/routes"
	"log"

	"github.com/spf13/viper"
)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("error trying read from config: %v", err)
	}

	if err := grpc.CreateAllConnections(); err != nil {
		log.Fatalf("error trying create all grpc clients: %v", err)
	}

	routes.GatewayApi()
}
