package controllers

import (
	"amaksimov/engineService/db"

	"amaksimov/pkg/grpc/pb"
	"amaksimov/pkg/grpc/statusCodes"
	"context"
	"database/sql"
	"errors"
	"log"

	"google.golang.org/grpc/status"
)

type EngineServer struct {
	pb.UnimplementedEngineServiceServer
}

func (es *EngineServer) GetAllEngines(ctx context.Context, req *pb.Req) (*pb.Engines, error) {
	engines, err := db.GetAllEngines()
	if err != nil {
		log.Printf("error querying all engines: %v", err)
		return nil, err
	}

	var enginesPB pb.Engines
	for _, val := range engines {
		enginesPB.Engines = append(enginesPB.Engines, &pb.Engine{Id: val.ID, Volume: val.Volume})
	}

	return &pb.Engines{Engines: enginesPB.Engines}, nil
}

func (es *EngineServer) GetEngineByID(ctx context.Context, req *pb.EngineID) (*pb.Engine, error) {
	engine, err := db.GetEngineByID(req.Id)
	if err != nil {
		log.Printf("error querying engine by id: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(statusCodes.NoData, "no rows")
		}

		return nil, err
	}

	return &pb.Engine{Id: engine.ID, Volume: engine.Volume}, nil
}

func (es *EngineServer) GetEnginesByIDs(ctx context.Context, req *pb.EnginesIDs) (*pb.Engines, error) {
	engines, err := db.GetEnginesByIDs(req.EnginesIDs)
	if err != nil {
		return nil, err
	}

	var enginesPB pb.Engines
	for _, val := range engines {
		enginesPB.Engines = append(enginesPB.Engines, &pb.Engine{Id: val.ID, Volume: val.Volume})
	}

	return &pb.Engines{Engines: enginesPB.Engines}, nil
}
