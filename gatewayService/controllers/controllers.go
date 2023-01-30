package controllers

import (
	"amaksimov/gatewayService/grpc"
	"amaksimov/gatewayService/models"

	"amaksimov/pkg/grpc/pb"
	"amaksimov/pkg/grpc/statusCodes"

	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/status"
)

func GetAllEngines(c echo.Context) error {
	engineService := pb.NewEngineServiceClient(grpc.EngineClient.Conn)

	response, err := engineService.GetAllEngines(context.Background(), &pb.Req{})
	if err != nil {
		log.Printf("error calling getAllEngines func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	return c.JSON(http.StatusOK, models.EnginesResponse{Data: response.Engines})
}

func GetUserCars(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID < 1 {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "invalid id"})
	}

	userService := pb.NewUserServiceClient(grpc.UserClient.Conn)

	userClientResp, err := userService.GetCarsIDsByID(context.Background(), &pb.UserID{Id: int32(userID)})
	if err != nil {
		log.Printf("error calling GetCarsIDsByID func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	carService := pb.NewCarServiceClient(grpc.CarClient.Conn)

	carClientResp, err := carService.GetCarsByIDs(context.Background(), &pb.CarsIDs{CarsIDs: userClientResp.CarsIDs})
	if err != nil {
		log.Printf("error calling GetCarsByIDs func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	var enginesIDs []int32
	for _, val := range carClientResp.Cars {
		enginesIDs = append(enginesIDs, val.Id)
	}

	engineService := pb.NewEngineServiceClient(grpc.EngineClient.Conn)

	engineClientResp, err := engineService.GetEnginesByIDs(context.Background(), &pb.EnginesIDs{EnginesIDs: enginesIDs})
	if err != nil {
		log.Printf("error calling GetEnginesByIDs func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	var resp []*models.Car

	for in, val := range carClientResp.Cars {
		var r = &models.Car{ID: val.Id,
			Concern: val.Concern,
			Model:   val.Model,
			Year:    val.Year,
			Engine: models.Engine{
				EngineID: engineClientResp.Engines[in].Id,
				Volume:   engineClientResp.Engines[in].Volume},
			Used:  val.Used,
			Price: val.Price}
		resp = append(resp, r)
	}

	return c.JSON(http.StatusOK, models.CarResponse{Data: resp})
}

func GetCarEngine(c echo.Context) error {
	carID, err := strconv.Atoi(c.Param("id"))
	if err != nil || carID < 1 {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "invalid id"})
	}

	carService := pb.NewCarServiceClient(grpc.CarClient.Conn)

	carClientResp, err := carService.GetEngineID(context.Background(), &pb.CarID{Id: int32(carID)})
	if status.Code(err) == statusCodes.NoData {
		return c.JSON(200, models.CarResponse{Data: nil})
	}
	if err != nil {
		log.Printf("error calling GetEngineID func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	engineService := pb.NewEngineServiceClient(grpc.EngineClient.Conn)

	engineClientResp, err := engineService.GetEngineByID(context.Background(), &pb.EngineID{Id: carClientResp.Id})
	if status.Code(err) == statusCodes.NoData {
		return c.JSON(200, models.EngineResponse{Data: nil})
	}
	if err != nil {
		log.Printf("error calling GetEnginebyIDfunc: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	return c.JSON(http.StatusOK, models.EngineResponse{Data: engineClientResp})
}

func GetUserEngines(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil || userID < 1 {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "invalid id"})
	}

	userService := pb.NewUserServiceClient(grpc.UserClient.Conn)

	userClientResp, err := userService.GetCarsIDsByID(context.Background(), &pb.UserID{Id: int32(userID)})
	if err != nil {
		log.Printf("error calling GetCarsIDs func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	carService := pb.NewCarServiceClient(grpc.CarClient.Conn)

	carClientResp, err := carService.GetEnginesIDs(context.Background(), &pb.CarsIDs{CarsIDs: userClientResp.CarsIDs})
	if err != nil {
		log.Printf("error calling GetEnginesIDsByIDS func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	engineService := pb.NewEngineServiceClient(grpc.EngineClient.Conn)

	engineClientResp, err := engineService.GetEnginesByIDs(context.Background(), &pb.EnginesIDs{EnginesIDs: carClientResp.EnginesIDs})
	if err != nil {
		log.Printf("error calling GetEnginesByIDs func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	return c.JSON(http.StatusOK, models.EnginesResponse{Data: engineClientResp.Engines})
}

func CreateEngine(c echo.Context) error {
	var e models.Engine

	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "invalid engine"})
	}

	engineService := pb.NewEngineServiceClient(grpc.EngineClient.Conn)

	engineServiceResp, err := engineService.CreateEngine(context.Background(), &pb.Engine{Volume: e.Volume})
	if err != nil {
		log.Printf("error calling CreateEngine func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	var en = &pb.Engine{Id: engineServiceResp.Id, Volume: engineServiceResp.Volume}

	return c.JSON(http.StatusCreated, models.EngineResponse{Data: en})
}

func UpdateEngine(c echo.Context) error {
	engineID, err := strconv.Atoi(c.Param("id"))
	if err != nil || engineID < 1 {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "invalid id"})
	}

	var e models.Engine

	if err := c.Bind(&e); err != nil {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "invalid engine"})
	}

	engineService := pb.NewEngineServiceClient(grpc.EngineClient.Conn)

	en, err := engineService.UpdateEngine(context.Background(), &pb.UpdateEngineRequest{Id: int32(engineID), Volume: e.Volume})
	if status.Code(err) == statusCodes.NoData {
		return c.JSON(200, models.EngineResponse{Data: nil})
	}
	if err != nil {
		log.Printf("error calling DeleteEngine func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	return c.JSON(http.StatusOK, models.EngineResponse{Data: en})
}

func DeleteEngine(c echo.Context) error {
	engineID, err := strconv.Atoi(c.Param("id"))
	if err != nil || engineID < 1 {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "invalid id"})
	}

	engineService := pb.NewEngineServiceClient(grpc.EngineClient.Conn)

	if _, err := engineService.DeleteEngine(context.Background(), &pb.EngineID{Id: int32(engineID)}); status.Code(err) == statusCodes.NoData {
		return c.JSON(200, models.EngineResponse{Data: nil})
	}
	if err != nil {
		log.Printf("error calling DeleteEngine func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	carService := pb.NewCarServiceClient(grpc.CarClient.Conn)

	if _, err := carService.DeleteCarConfiguration(context.Background(), &pb.EngineID{Id: int32(engineID)}); status.Code(err) == statusCodes.NoData {
		return c.JSON(200, models.EngineResponse{Data: nil})
	}
	if err != nil {
		log.Printf("error calling DeleteEngine func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	return c.JSON(http.StatusOK, models.Response{Message: "Successfully deleted"})
}

func GetEngineByID(c echo.Context) error {
	engineID, err := strconv.Atoi(c.Param("id"))
	if err != nil || engineID < 1 {
		return c.JSON(http.StatusBadRequest, models.Response{Message: "invalid id"})
	}

	engineService := pb.NewEngineServiceClient(grpc.EngineClient.Conn)

	engineResponse, err := engineService.GetEngineByID(context.Background(), &pb.EngineID{Id: int32(engineID)})
	if status.Code(err) == statusCodes.NoData {
		return c.JSON(200, models.EngineResponse{Data: nil})
	}
	if err != nil {
		log.Printf("error calling GetEngineByID func: %v", err)
		return c.JSON(http.StatusInternalServerError, models.Response{Message: "something went wrong"})
	}

	return c.JSON(http.StatusOK, models.EngineResponse{Data: engineResponse})
}
