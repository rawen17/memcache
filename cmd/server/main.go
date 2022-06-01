package main

import (
	"log"
	"net"

	"memcache/pkg/proto/person"
	svc "memcache/pkg/service"
	persongrpc "memcache/pkg/service/transport/grpc"
	"memcache/pkg/storage/mystorage"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)

	store := mystorage.NewMyStorage()
	service, err := svc.NewService(store)
	if err != nil {
		log.Fatalf("Failed to create service: %v", err)
	}
	srv, err := persongrpc.NewServer(service)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	person.RegisterPersonServiceServer(s, srv)

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
