package main

import (
	"net"
	"os"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func main() {
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
	err = grpcServer.Serve(listener)
	if err != nil {
		logger.Println(err.Error())
	}
}
