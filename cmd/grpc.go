package main

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func initRPC() *grpc.Server {
	s := grpc.NewServer()
	reflection.Register(s)
	return s
}
