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

	keys := cache.Redis.Client.Keys(ctx, "*engine_*") // TODO: что это за префикс/суффикс ключей? где обработка ошибки?
	if len(keys.Val()) != 0 {
		for _, key := range keys.Val() {
			value := cache.Redis.Client.Get(ctx, key).Val() // TODO: где обработка ошибки?

			if value == "" {
				log.Println("cant find record in cache")

				continue
			}

			var e models.Engine

			if err := json.Unmarshal([]byte(value), &e); err != nil {
				log.Printf("error trying unmarshall key bytes: %v", err)

				continue
			}

			engines = append(engines, e)
		}
	}

	if len(engines) != 0 { // TODO: ну взял ты одну запись из redis, где тут All ?
		goto Parse
	}

	engines, err = db.GetAllEngines()
	if err != nil {
		log.Printf("error querying all engines: %v", err)

		return nil, err
	}
Parse:
	var enginesPB pb.Engines // TODO: нет капасити для enginesPB.Engines, хотя размер известен

	for _, val := range engines { // TODO: используй индексы
		enginesPB.Engines = append(enginesPB.Engines, &pb.Engine{Id: val.ID, Volume: val.Volume})
	}

	return &pb.Engines{Engines: enginesPB.Engines}, nil
}

func (es *EngineServer) GetEngineByID(ctx context.Context, req *pb.EngineID) (*pb.Engine, error) {
	var (
		engine *models.Engine
		err    error
	)

	key := fmt.Sprintf("engine_%v", req.Id) // TODO: req.Id какой тип? - почитай за форматирование и почему тут не нужен %v

	value := cache.Redis.Client.Get(ctx, key).Val() // TODO: где обработка ошибки?

	if value == "" {
		log.Println("cant find record in cache")
	}

	// TODO: что ты тут пытаешься декодировать, если у тебя value пустой?
	if err := json.Unmarshal([]byte(value), engine); err != nil {
		log.Printf("error trying unmarshall key bytes: %v", err)
	}

	// TODO: хочу видеть реализацию без goto
	if engine != nil {
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
		key := fmt.Sprintf("engine_%d", val) // TODO: зачем везде писать engine_, почему не вынес в константу?

		value := cache.Redis.Client.Get(ctx, key).Val() // TODO: где обработка ошибки?

		if value == "" {
			log.Println("cant find record in cache")

			continue
		}

		var e models.Engine

		if err := json.Unmarshal([]byte(value), &e); err != nil {
			log.Printf("error trying unmarshall key bytes: %v", err)

			continue
		}

		engines = append(engines, e)
	}

	// TODO: убрать везде goto, и написать логику без него, так же, почему в redis записывается только при запросе одного двигателя,
	// а при ALL / ByIDs нет?
	if len(engines) != 0 {
		goto Parse
	}

	engines, err = db.GetEnginesByIDs(req.EnginesIDs)
	if err != nil {
		return nil, err
	}

Parse:
	var enginesPB pb.Engines // TODO: тот же комментарий про размер, он заранее известен
	for _, val := range engines {
		enginesPB.Engines = append(enginesPB.Engines, &pb.Engine{Id: val.ID, Volume: val.Volume})
	}

	return &pb.Engines{Engines: enginesPB.Engines}, nil
}

func (es *EngineServer) CreateEngine(ctx context.Context, req *pb.Engine) (*pb.Engine, error) {
	engine, err := db.CreateEngine(models.Engine{ID: req.Id, Volume: req.Volume})
	if err != nil {
		return nil, err // TODO: где враппинг ошибки?
	}

	var e = models.Engine{ID: engine.ID, Volume: engine.Volume} // TODO: обзывай переменные нормально, 'e' выглядит больше как псевдоним объекта для которого расписывается метод

	engineBytes, err := json.Marshal(&e)
	if err != nil {
		return nil, err // TODO: где враппинг ошибки?
	}

	key := fmt.Sprintf("engine_%d", engine.ID)

	cache.Redis.Client.Set(context.Background(), key, engineBytes, 0) // TODO: где обработка ошибки?

	return &pb.Engine{Id: e.ID, Volume: e.Volume}, nil
}

// TODO: все отрефакторить согласно предыдущим комментариям
func (es *EngineServer) UpdateEngine(ctx context.Context, req *pb.UpdateEngineRequest) (*pb.Engine, error) {
	engine, err := db.UpdateEngine(models.Engine{ID: req.Id, Volume: req.Volume})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(statusCodes.NoData, "no rows")
		}
		return nil, err
	}

	var e = models.Engine{ID: engine.ID, Volume: engine.Volume}

	engineBytes, err := json.Marshal(&e)
	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("engine_%d", req.Id)

	cache.Redis.Client.Set(context.Background(), key, engineBytes, 0)

	return &pb.Engine{Id: e.ID, Volume: e.Volume}, nil
}

// TODO: все отрефакторить согласно предыдущим комментариям
func (es *EngineServer) DeleteEngine(ctx context.Context, req *pb.EngineID) (*pb.Resp, error) {
	if err := db.DeleteEngine(req.Id); err != nil {
		return nil, err
	}

	key := fmt.Sprintf("engine_%d", req.Id)

	cache.Redis.Client.Del(context.Background(), key)

	return &pb.Resp{}, nil
}
