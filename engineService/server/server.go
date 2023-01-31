package server

import (
	"amaksimov/engineService/controllers"
	"amaksimov/pkg/grpc/pb"
	"fmt"

	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func CreateEngineServer() error {
	address := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))

	listener, err := net.Listen(viper.GetString("server.network"), address)
	if err != nil {
		return err // TODO: где враппинг ошибки?
	}

	grpcServer := grpc.NewServer()

	pb.RegisterEngineServiceServer(grpcServer, &controllers.EngineServer{})

	log.Printf("grpc server listening on: %v", viper.GetString("server.port"))

	if err := grpcServer.Serve(listener); err != nil {
		return err // TODO: где враппинг ошибки?
	}

	return nil
}
