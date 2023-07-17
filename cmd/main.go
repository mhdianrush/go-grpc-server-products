package main

import (
	"go-gRPC-server-products/cmd/config"
	"go-gRPC-server-products/cmd/services"
	productpb "go-gRPC-server-products/pb/product"
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
	db := config.ConnectDB()

	logger := logrus.New()

	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logger.Println(err.Error())
	}
	logger.SetOutput(file)

	logger.Println("Server Running on Port 8080")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		logger.Println(err.Error())
	}

	grpcServer := grpc.NewServer()

	productService := services.ProductService{DB: db}
	productpb.RegisterProductServiceServer(grpcServer, &productService)

	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Println(err.Error())
	}
}
