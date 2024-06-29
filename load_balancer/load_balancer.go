package main

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

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
	start := time.Now()
	log.Printf("[%s] Load Balancer received Get request, forwarding to server on port %s", start.Format("2006-01-02 15:04:05.000000"), port)
	res, err := client.Get(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] Second Get Response from Server %s: Data: %v", time.Now().Format("2006-01-02 15:04:05.000000"), res.ServerPort, res.Data)
	return res, nil
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	client, port := s.getNextServer()
	start := time.Now()
	log.Printf("[%s] Load Balancer received Set request, forwarding to server on port %s", start.Format("2006-01-02 15:04:05.000000"), port)
	res, err := client.Set(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] Set Response from Server %s", time.Now().Format("2006-01-02 15:04:05.000000"), port)
	return res, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	client, port := s.getNextServer()
	start := time.Now()
	log.Printf("[%s] Load Balancer received Delete request, forwarding to server on port %s", start.Format("2006-01-02 15:04:05.000000"), port)
	res, err := client.Delete(ctx, req)
	if err != nil {
		return nil, err
	}
	log.Printf("[%s] Delete Response from Server %s", time.Now().Format("2006-01-02 15:04:05.000000"), port)
	return res, nil
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
