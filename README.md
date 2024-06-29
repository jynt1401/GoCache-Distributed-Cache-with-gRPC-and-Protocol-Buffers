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
