package main

import (
	"context"
	"log"
	"strconv"
	"sync"
	"time"

	pb "gocache/proto"

	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:50050", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()
	client := pb.NewCacheServiceClient(conn)

	var wgSet, wgGet1, wgDelete, wgGet2 sync.WaitGroup

	// Concurrent Set requests
	for i := 0; i < 3; i++ {
		wgSet.Add(1)
		go func(index int) {
			defer wgSet.Done()
			key := "key" + strconv.Itoa(index+1)
			value := "value" + strconv.Itoa(index+1)
			_, err := client.Set(context.Background(), &pb.SetRequest{Key: key, Value: value})
			if err != nil {
				log.Fatalf("could not set: %v", err)
			}
			log.Printf("Set Request for %s=%s", key, value)
		}(i)
	}
	wgSet.Wait()
	time.Sleep(2 * time.Second) // Timeout after Set requests

	// Concurrent Get requests after Set operations
	for i := 0; i < 10; i++ {
		wgGet1.Add(1)
		go func(index int) {
			defer wgGet1.Done()
			time.Sleep(time.Millisecond * 20) // Simulating delay after Set operations
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			res, err := client.Get(ctx, &pb.GetRequest{})
			if err != nil {
				log.Fatalf("could not get: %v", err)
			}
			log.Printf("Get Response %d from Server %s at %s: Data: %v", index+1, res.ServerPort, start.Format("15:04:05.000000000"), res.Data)
		}(i)
	}
	wgGet1.Wait()
	time.Sleep(2 * time.Second) // Timeout after first Get requests

	// Concurrent Delete requests
	for i := 0; i < 2; i++ {
		wgDelete.Add(1)
		go func(index int) {
			defer wgDelete.Done()
			time.Sleep(time.Millisecond * 30) // Simulating delay between Delete operations
			key := "key" + strconv.Itoa(index+1)
			_, err := client.Delete(context.Background(), &pb.DeleteRequest{Key: key})
			if err != nil {
				log.Fatalf("could not delete: %v", err)
			}
			log.Printf("Delete Request for %s", key)
		}(i)
	}
	wgDelete.Wait()
	time.Sleep(2 * time.Second) // Timeout after Delete requests

	// Concurrent Get requests after Delete operations
	for i := 0; i < 10; i++ {
		wgGet2.Add(1)
		go func(index int) {
			defer wgGet2.Done()
			time.Sleep(time.Millisecond * 40) // Simulating delay after Delete operations
			start := time.Now()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			res, err := client.Get(ctx, &pb.GetRequest{})
			if err != nil {
				log.Fatalf("could not get: %v", err)
			}
			log.Printf("Second Get Response %d from Server %s at %s: Data: %v", index+1, res.ServerPort, start.Format("15:04:05.000000000"), res.Data)
		}(i)
	}
	wgGet2.Wait()
}
