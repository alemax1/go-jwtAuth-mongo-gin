package server

import (
	"amaksimov/engineService/controllers"
	"amaksimov/pkg/grpc/pb"
	"log"
	"net"

	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

func CreateEngineServer() error {
	port := viper.GetString("server.port")
	lis, err := net.Listen(viper.GetString("server.network"), ":"+port)
	if err != nil {
		return err
	}

	grpcServer := grpc.NewServer()

	pb.RegisterEngineServiceServer(grpcServer, &controllers.EngineServer{})

	log.Printf("grpc server listening on: %v", port)
	if err := grpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
