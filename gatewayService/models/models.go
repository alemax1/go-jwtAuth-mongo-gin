package models

import (
	"amaksimov/pkg/grpc/pb"
	_ "encoding/json"
)

type Response struct {
	Message string `json:"message"`
}

type Car struct {
	ID      int32  `json:"ID"`
	Concern string `json:"concern"`
	Model   string `jsom:"model"`
	Year    int32  `json:"year"`
	Used    bool   `json:"used"`
	Engine  `json:"engine"`
	Price   int32 `json:"price"`
}

type CarResponse struct {
	Data []*Car `json:"data"`
}

type Engine struct {
	EngineID int32   `json:"ID"`
	Volume   float32 `json:"volume"`
}

type EngineResponse struct {
	Data *pb.Engine `json:"data"`
}

type EnginesResponse struct {
	Data []*pb.Engine `json:"data"`
}
