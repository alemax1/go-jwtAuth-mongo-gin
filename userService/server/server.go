package server

import (
	"amaksimov/pkg/grpc/pb"
	"amaksimov/userService/controllers"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func CreateUserServer() error {
	port := viper.GetString("server.port")

	lis, err := net.Listen(viper.GetString("server.network"), ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &controllers.UserServer{})

	log.Printf("grpc server listening on: %v", port)

	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
