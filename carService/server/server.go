package server

import (
	"amaksimov/carService/controllers"
	"amaksimov/pkg/grpc/pb"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func CreateCarServer() error {
	port := viper.GetString("server.port")

	// TODO: почему только порт, почему на задается ip ?
	lis, err := net.Listen(viper.GetString("server.network"), ":"+port)
	if err != nil {
		return err // TODO: где враппинг ошибки?
	}

	grpcServer := grpc.NewServer()

	pb.RegisterCarServiceServer(grpcServer, &controllers.CarServer{})

	log.Printf("grpc server listening on: %v", port)

	if err := grpcServer.Serve(lis); err != nil {
		return err // TODO: где враппинг ошибки?
	}

	return nil
}
