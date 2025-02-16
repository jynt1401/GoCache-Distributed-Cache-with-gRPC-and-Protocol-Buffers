# GoCache: Distributed Cache with gRPC and Protocol Buffers

GoCache is a distributed caching system built using Go (Golang), gRPC, and Protocol Buffers. It allows multiple servers to store and retrieve key-value pairs efficiently, ensuring data consistency and high availability.

## Features

- **Distributed Caching:** Utilizes multiple servers to store cached data, enabling efficient data access across a network.
- **gRPC Communication:** Uses gRPC for fast and efficient communication between clients and servers, utilizing Protocol Buffers for defining service interfaces.
- **Data Replication:** Implements replication mechanisms to synchronize data changes across all cache servers, ensuring data consistency.
- **Load Balancing:** Includes a load balancer to evenly distribute client requests among multiple cache servers, optimizing system performance.

## How to Run

### Prerequisites

- Go (Golang) installed on your machine.
- Protocol Buffers compiler (`protoc`) installed.
- Basic understanding of terminal commands.

### Running the Code

1. **Clone the Repository:**

   ```bash
   git clone https://github.com/jynt1401/GoCache-Distributed-Cache-with-gRPC-and-Protocol-Buffers.git
   cd GoCache
   ```
2. **Generate .proto file**

   ```bash
   protoc --go_out=. --go-grpc_out=. ./proto/cache.proto
   ```
3. **Run Servers**

   ```bash
   cd .\server\
   ```
   open 5 terminal and run
    ```bash
   go run server.go -port=50051
   ```
    ```bash
   go run server.go -port=50052
   ```
    ```bash
   go run server.go -port=50053
   ```
    ```bash
   go run server.go -port=50054
   ```
    ```bash
   go run server.go -port=50055
   ```

4. **Run Load Balancer**

   ```bash
   cd .\load_balancer\
   go run load_balancer.go
   ```

5. **Run Client**

   ```bash
   cd .\client\
   go run client.go
   ```
