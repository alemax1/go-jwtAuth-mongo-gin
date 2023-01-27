package controllers

import (
	"amaksimov/pkg/grpc/pb"
	"amaksimov/userService/db"
	"context"
	"log"
)

type UserServer struct {
	pb.UnimplementedUserServiceServer
}

func (us *UserServer) GetCarsIDsByID(ctx context.Context, req *pb.UserID) (*pb.CarsIDs, error) {
	carsIDs, err := db.GetCarsIDsByID(req.Id)
	if err != nil {
		log.Printf("error querying cars ids by id: %v", err)
		return nil, err
	}

	return &pb.CarsIDs{CarsIDs: carsIDs}, nil
}
