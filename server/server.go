package main

import (
	"context"
	"flag"
	"log"
	"net"
	"sync"

	pb "gocache/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCacheServiceServer
	data map[string]string
	port string
	mu   sync.Mutex
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[req.Key] = req.Value
	return &pb.SetResponse{Message: "Set operation successful"}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return &pb.GetResponse{Data: s.data, ServerPort: s.port}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, req.Key)
	return &pb.DeleteResponse{Message: "Delete operation successful"}, nil
}

func main() {
	port := flag.String("port", "50051", "Port on which the server will listen")
	flag.Parse()

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *port, err)
	}
	s := grpc.NewServer()
	pb.RegisterCacheServiceServer(s, &server{data: make(map[string]string), port: *port})
	log.Printf("Server listening on port %s", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
