package grpc

import (
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type grpcClient struct {
	Conn *grpc.ClientConn
}

var (
	EngineClient grpcClient
	UserClient   grpcClient
	CarClient    grpcClient
)

func newGRPCConnection(address string) (*grpc.ClientConn, error) {
	var err error

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func CreateAllConnections() error {
	var err error

	EngineClient.Conn, err = newGRPCConnection(viper.GetString("servicesURLs.engine"))
	if err != nil {
		return err
	}

	UserClient.Conn, err = newGRPCConnection(viper.GetString("servicesURLs.user"))
	if err != nil {
		return err
	}

	CarClient.Conn, err = newGRPCConnection(viper.GetString("servicesURLs.car"))
	if err != nil {
		return err
	}

	return nil
}
