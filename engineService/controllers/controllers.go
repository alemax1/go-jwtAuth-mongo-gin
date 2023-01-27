package controllers

import (
	"amaksimov/engineService/cache"
	"amaksimov/engineService/db"
	"amaksimov/engineService/models"
	"encoding/json"
	"fmt"

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
	var (
		engines []models.Engine
		err     error
	)

	keys := cache.Redis.Client.Keys(ctx, "*engine_*")
	if len(keys.Val()) != 0 {
		for _, val := range keys.Val() {
			key := cache.Redis.Client.Get(ctx, val)

			var e models.Engine

			if err := json.Unmarshal([]byte(key.Val()), &e); err != nil {
				log.Printf("error trying unmarshall key bytes: %v", err)

				continue
			}

			engines = append(engines, e)
		}
		goto Parse
	}

	engines, err = db.GetAllEngines()
	if err != nil {
		log.Printf("error querying all engines: %v", err)

		return nil, err
	}
Parse:
	var enginesPB pb.Engines

	for _, val := range engines {
		enginesPB.Engines = append(enginesPB.Engines, &pb.Engine{Id: val.ID, Volume: val.Volume})
	}

	return &pb.Engines{Engines: enginesPB.Engines}, nil
}

func (es *EngineServer) GetEngineByID(ctx context.Context, req *pb.EngineID) (*pb.Engine, error) {
	var (
		engine *models.Engine
		err    error
	)

	pattern := fmt.Sprintf("engines_%v", req.Id)

	keys := cache.Redis.Client.Keys(ctx, pattern)
	if len(keys.Val()) != 0 {
		key := cache.Redis.Client.Get(ctx, keys.Val()[0])

		if err := json.Unmarshal([]byte(key.Val()), &engine); err != nil {
			log.Printf("error trying unmarshall key bytes: %v", err)
		}
		goto Return
	}

	engine, err = db.GetEngineByID(req.Id)
	if err != nil {
		log.Printf("error querying engine by id: %v", err)
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(statusCodes.NoData, "no rows")
		}

		return nil, err
	}
Return:
	return &pb.Engine{Id: engine.ID, Volume: engine.Volume}, nil
}

func (es *EngineServer) GetEnginesByIDs(ctx context.Context, req *pb.EnginesIDs) (*pb.Engines, error) {
	// var (
	// 	engines []models.Engine
	// 	err     error
	// 	IDs     string
	// )

	// for _, val := range req.EnginesIDs {
	// 	IDs = IDs + strconv.Itoa(int(val))
	// }

	// fmt.Println(IDs)
	// pattern := fmt.Sprintf("engines_[%v]", IDs)
	// fmt.Println(pattern)

	// keys := cache.Redis.Client.Keys(ctx, pattern)
	// fmt.Println(keys)
	// if len(keys.Val()) != 0 {
	// 	for _, val := range keys.Val() {
	// 		key := cache.Redis.Client.Get(ctx, val)

	// 		var e models.Engine

	// 		if err := json.Unmarshal([]byte(key.Val()), &e); err != nil {
	// 			log.Printf("error trying unmarshall key bytes: %v", err)

	// 			continue
	// 		}

	// 		engines = append(engines, e)
	// 	}
	// 	fmt.Println("CACHEEEEEEEEEEE")
	// 	goto Parse
	// }

	engines, err := db.GetEnginesByIDs(req.EnginesIDs)
	if err != nil {
		return nil, err
	}

	// Parse:
	var enginesPB pb.Engines
	for _, val := range engines {
		enginesPB.Engines = append(enginesPB.Engines, &pb.Engine{Id: val.ID, Volume: val.Volume})
	}

	return &pb.Engines{Engines: enginesPB.Engines}, nil
}
