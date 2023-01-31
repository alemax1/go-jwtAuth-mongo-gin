package server

import (
	"amaksimov/pkg/grpc/pb"
	"amaksimov/userService/controllers"
	"fmt"
	"log"
	"net"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func CreateUserServer() error {
	address := fmt.Sprintf("%s:%s", viper.GetString("server.host"), viper.GetString("server.port"))

	listener, err := net.Listen(viper.GetString("server.network"), address)
	if err != nil {
		return errors.Wrap(err, "listen failed")
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &controllers.UserServer{})

	log.Printf("grpc server listening on: %v", viper.GetString("server.port"))

	if err := grpcServer.Serve(listener); err != nil {
		return errors.Wrap(err, "failed creating grpc server")
	}

	return nil
}
