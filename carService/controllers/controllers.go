package controllers

import (
	"amaksimov/carService/db"

	"amaksimov/pkg/grpc/pb"
	"amaksimov/pkg/grpc/statusCodes"
	"context"
	"database/sql"
	"errors"
	"log"

	"google.golang.org/grpc/status"
)

type CarServer struct {
	pb.UnimplementedCarServiceServer
}

func (cs *CarServer) GetCarsByIDs(ctx context.Context, req *pb.CarsIDs) (*pb.Cars, error) {
	cars, err := db.GetCarsByIDs(req.CarsIDs)
	if err != nil {
		log.Printf("error querying cars by IDs: %v", err)
		return nil, err
	}

	var carsPB pb.Cars
	for _, val := range cars {
		carsPB.Cars = append(carsPB.Cars, &pb.Car{Id: val.ID,
			Concern:  val.Concern,
			Model:    val.Model,
			Year:     val.Year,
			Used:     val.Used,
			EngineID: val.EngineID,
			Price:    val.Price})
	}

	return &pb.Cars{Cars: carsPB.Cars}, nil
}

func (cs *CarServer) GetEngineID(ctx context.Context, req *pb.CarID) (*pb.EngineID, error) {
	engineID, err := db.GetEngineID(req.Id)
	if err != nil {
		log.Printf("error querying engines by ID: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(statusCodes.NoData, "no rows")
		}

		return nil, err
	}

	return &pb.EngineID{Id: engineID}, nil
}

func (cs *CarServer) GetEnginesIDs(ctx context.Context, req *pb.CarsIDs) (*pb.EnginesIDs, error) {
	enginesIDs, err := db.GetEnginesByIDs(req.CarsIDs)
	if err != nil {
		log.Printf("error querying engines by IDs: %v", err)
		return nil, err
	}

	return &pb.EnginesIDs{EnginesIDs: enginesIDs}, nil
}

func (cs *CarServer) DeleteCarConfiguration(ctx context.Context, req *pb.EngineID) (*pb.Resp, error) {
	if err := db.DeleteCarConfiguration(req.Id); err != nil {
		log.Printf("error trying delete car configuration: %v", err)
		return nil, err
	}

	return &pb.Resp{}, nil
}
