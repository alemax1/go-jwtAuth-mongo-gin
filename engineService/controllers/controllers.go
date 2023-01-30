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
	}

	if len(engines) != 0 {
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

	pattern := fmt.Sprintf("engine_%v", req.Id)

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
	var (
		engines []models.Engine
		err     error
	)

	for _, val := range req.EnginesIDs {
		key := fmt.Sprintf("engine_%d", val)

		value := cache.Redis.Client.Get(ctx, key).Val()

		var e models.Engine

		if err := json.Unmarshal([]byte(value), &e); err != nil {
			log.Printf("error trying unmarshall key bytes: %v", err)

			continue
		}

		engines = append(engines, e)
	}

	if len(engines) != 0 {
		goto Parse
	}

	engines, err = db.GetEnginesByIDs(req.EnginesIDs)
	if err != nil {
		return nil, err
	}

Parse:
	var enginesPB pb.Engines
	for _, val := range engines {
		enginesPB.Engines = append(enginesPB.Engines, &pb.Engine{Id: val.ID, Volume: val.Volume})
	}

	return &pb.Engines{Engines: enginesPB.Engines}, nil
}

func (es *EngineServer) CreateEngine(ctx context.Context, req *pb.Engine) (*pb.Engine, error) {
	engine, err := db.CreateEngine(req)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("engine_%d", engine.ID)

	var e = &models.Engine{ID: engine.ID, Volume: engine.Volume}

	engineBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	cache.Redis.Client.Set(context.Background(), key, engineBytes, 0)

	return &pb.Engine{Id: e.ID, Volume: e.Volume}, nil
}

func (es *EngineServer) UpdateEngine(ctx context.Context, req *pb.UpdateEngineRequest) (*pb.Engine, error) {
	engine, err := db.UpdateEngine(req.Id, req.Volume)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("engine_%d", req.Id)

	var e = &models.Engine{ID: engine.ID, Volume: engine.Volume}

	engineBytes, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	cache.Redis.Client.Set(context.Background(), key, engineBytes, 0)

	fmt.Println(e)

	return &pb.Engine{Id: e.ID, Volume: e.Volume}, nil
}

func (es *EngineServer) DeleteEngine(ctx context.Context, req *pb.EngineID) (*pb.Resp, error) {
	if err := db.DeleteEngine(req.Id); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("engine_%d", req.Id)

	cache.Redis.Client.Del(context.Background(), key)

	return &pb.Resp{}, nil
}
