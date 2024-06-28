package main

import (
	"context"
	"log"
	"net"
	"sync"

	pb "gocache/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCacheServiceServer
	clients []pb.CacheServiceClient
	ports   []string
	index   int
	mu      sync.Mutex
}

func (s *server) getNextServer() (pb.CacheServiceClient, string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	client := s.clients[s.index]
	port := s.ports[s.index]
	s.index = (s.index + 1) % len(s.clients)
	return client, port
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	client, port := s.getNextServer()
	log.Printf("Load Balancer received Get request, forwarding to server on port %s", port)
	return client.Get(ctx, req)
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	client, port := s.getNextServer()
	log.Printf("Load Balancer received Set request, forwarding to server on port %s", port)
	return client.Set(ctx, req)
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	client, port := s.getNextServer()
	log.Printf("Load Balancer received Delete request, forwarding to server on port %s", port)
	return client.Delete(ctx, req)
}

func main() {
	ports := []string{":50051", ":50052", ":50053", ":50054", ":50055"} // Ports for servers
	var clients []pb.CacheServiceClient

	for _, port := range ports {
		conn, err := grpc.Dial("localhost"+port, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server: %v", err)
		}
		clients = append(clients, pb.NewCacheServiceClient(conn))
	}

	lis, err := net.Listen("tcp", ":50050")
	if err != nil {
		log.Fatalf("Failed to listen on port 50050: %v", err)
	}

	s := grpc.NewServer()
	pb.RegisterCacheServiceServer(s, &server{clients: clients, ports: ports})

	log.Printf("Load Balancer listening on port 50050")
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
