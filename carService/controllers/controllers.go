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

	var carsPB pb.Cars         // TODO: ну добавь ты разделение строки от объявления переменной и выполнения цикла, зачем в кучу все сводить?
	for _, val := range cars { // TODO: почему бы не использовать индекс вместо постоянного копирования значений в val ?
		carsPB.Cars = append(carsPB.Cars, &pb.Car{Id: val.ID, // TODO: почему carsPB.Cars не задать капасити, а уже после добавлять элементы?
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
	engineID, err := db.GetEngineID(req.Id) // TODO: что за Id ?
	if err != nil {
		log.Printf("error querying engines by ID: %v", err) // TODO: почему не разделить лог и проверку ошибки новой строкой?
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Error(statusCodes.NoData, "no rows")
		}

		return nil, err
	}

	return &pb.EngineID{Id: engineID}, nil // TODO: что за Id ?
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
	if err := db.DeleteCarConfiguration(req.Id); err != nil { // TODO: что за Id ? не раз было сказано ID
		log.Printf("error trying delete car configuration: %v", err)
		return nil, err
	}

	return &pb.Resp{}, nil // TODO: что за пустой response, почему так?
}
