package main

import (
	"context"
	"flag"
	"log"
	"net"
	"sync"
	"time"

	pb "gocache/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedCacheServiceServer
	data    map[string]string
	clients []pb.CacheServiceClient // Clients to other servers
	port    string
	mu      sync.Mutex
}

// ReplicateSet replicates a Set operation to all other servers
func (s *server) ReplicateSet(key, value string) {
	for _, client := range s.clients {
		go func(c pb.CacheServiceClient) {
			_, err := c.Set(context.Background(), &pb.SetRequest{Key: key, Value: value})
			if err != nil {
				log.Printf("Error replicating Set operation: %v", err)
			}
		}(client)
	}
}

// ReplicateDelete replicates a Delete operation to all other servers
func (s *server) ReplicateDelete(key string) {
	for _, client := range s.clients {
		go func(c pb.CacheServiceClient) {
			_, err := c.Delete(context.Background(), &pb.DeleteRequest{Key: key})
			if err != nil {
				log.Printf("Error replicating Delete operation: %v", err)
			}
		}(client)
	}
}

func (s *server) Set(ctx context.Context, req *pb.SetRequest) (*pb.SetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[req.Key] = req.Value

	// Replicate Set operation to other servers asynchronously
	go s.ReplicateSet(req.Key, req.Value)

	return &pb.SetResponse{Message: "Set operation successful"}, nil
}

func (s *server) Get(ctx context.Context, req *pb.GetRequest) (*pb.GetResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Simulate some processing time
	time.Sleep(time.Millisecond * 100)

	return &pb.GetResponse{Data: s.data, ServerPort: s.port}, nil
}

func (s *server) Delete(ctx context.Context, req *pb.DeleteRequest) (*pb.DeleteResponse, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.data, req.Key)

	// Replicate Delete operation to other servers asynchronously
	go s.ReplicateDelete(req.Key)

	return &pb.DeleteResponse{Message: "Delete operation successful"}, nil
}

func main() {
	port := flag.String("port", "50051", "Port on which the server will listen")
	flag.Parse()

	ports := []string{":50051", ":50052", ":50053", ":50054", ":50055"} // Ports for servers
	var clients []pb.CacheServiceClient

	for _, p := range ports {
		// Skip own port
		if p == ":"+*port {
			continue
		}
		conn, err := grpc.Dial("localhost"+p, grpc.WithInsecure())
		if err != nil {
			log.Fatalf("Failed to connect to server on port %s: %v", p, err)
		}
		clients = append(clients, pb.NewCacheServiceClient(conn))
	}

	lis, err := net.Listen("tcp", ":"+*port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", *port, err)
	}
	s := grpc.NewServer()
	pb.RegisterCacheServiceServer(s, &server{data: make(map[string]string), clients: clients, port: *port})
	log.Printf("Server listening on port %s", *port)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
