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

	// TOOD: что за название lis? - тот же комменатрий, почему без ip ?
	lis, err := net.Listen(viper.GetString("server.network"), ":"+port)
	if err != nil {
		return err // TODO: где враппинг ошибок? - тебе не раз за это говорили, в чем проблема с ними?
	}

	grpcServer := grpc.NewServer()

	pb.RegisterUserServiceServer(grpcServer, &controllers.UserServer{})

	log.Printf("grpc server listening on: %v", port)

	if err := grpcServer.Serve(lis); err != nil {
		return err // TODO: где враппинг ошибок?
	}

	return nil
}
