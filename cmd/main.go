package main

import (
	"go-gRPC-server-products/cmd/config"
	"go-gRPC-server-products/cmd/services"
	productpb "go-gRPC-server-products/pb/product"
	"net"
	"os"

	"github.com/joho/godotenv"
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

	if err = godotenv.Load(); err != nil {
		logger.Printf("failed load env file %s", err.Error())
	}

	listener, err := net.Listen("tcp", ":"+os.Getenv("SERVER_PORT"))
	if err != nil {
		logger.Printf("failed listen on port %s, %s", os.Getenv("SERVER_PORT"), err.Error())
	}

	logger.Printf("server running on port %s", os.Getenv("SERVER_PORT"))

	grpcServer := grpc.NewServer()

	productService := services.ProductService{DB: db}
	productpb.RegisterProductServiceServer(grpcServer, &productService)

	if err = grpcServer.Serve(listener); err != nil {
		logger.Printf("failed connect to server %s", err.Error())
	}
}
