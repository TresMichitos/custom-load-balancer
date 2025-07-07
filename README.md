[![tests](https://github.com/TresMichitos/custom-load-balancer/actions/workflows/tests.yml/badge.svg)](https://github.com/TresMichitos/custom-load-balancer/actions/workflows/tests.yml)
[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)

# Custom Load Balancer Using Go

Custom implementation of a load balancer using Go.

## Installation

1. Download the repository with:

   - ```
     git clone https://github.com/TresMichitos/custom-load-balancer.git
     ```
   - Or by downloading as a ZIP file

2. Enter the project directory with:
   ```
   cd custom-load-balancer/
   ```

## Usage

Run with:
  ```
  go run ./cmd/custom-load-balancer/ --algorithm {'RoundRobin'|...}
  ```
Or:
  ```
  docker compose up --build
  ```

Can send requests to load balancer with:
  ```
  curl http://localhost:8080
  ```
Or:
  ```
  go run demo/client/client.go http://localhost:8080 15
  ```

Shut down containers with:
  ```
  docker run compose down
  ```

---

## Development

### Testing:

Run Unit Tests with:
  ```
  go test ./internal/...
  ```

Run Integration Tests with:

1. Start Docker Containers with:
   ```
   sudo docker compose up
   ```
2. Run Integration Tests with:
   ```
   go test ./integration/...
   ```
3. Stop Docker Containers with: 
   ```
   sudo docker compose down
   ```

### Goals

- Use algorithm to route HTTP requests between multiple servers
- Make use of Go concurrency features (Goroutines, mutexes, channels etc)
- Host multiple web servers to be routed to via Docker Compose
- Record statistics on effect of load balancer on metrics such as latency, processing load etc
- Setup GitHub workflow testing
- Possibly do health checks on each server
- Possibly experiment with different load balancing algorithms

### Resources

- [A Tour of Go](https://go.dev/tour/list)
- [What is load balancing?](https://www.cloudflare.com/en-gb/learning/performance/what-is-load-balancing/)
- [Round Robin Load Balancing](https://www.vmware.com/topics/round-robin-load-balancing)
- [Example Go load balancer](https://dev.to/vivekalhat/building-a-simple-load-balancer-in-go-70d)
- [Docker Example](https://docs.docker.com/get-started/workshop/02_our_app/)
- [How Docker Compose works](https://docs.docker.com/compose/intro/compose-application-model/)

#### Docker Resources

- [Setup](https://medium.com/@aedemirsen/load-balancing-with-nginx-c1f19840e29)
- [Docker Containers](https://medium.com/@aedemirsen/load-balancing-with-docker-compose-and-nginx-b9077696f624)

### Possible Libraries To Use

- Included in Go standard library:
  - net/http
  - sync

---

## File Structure

```
.                                        # Root
├── cmd/
│   └── custom-load-balancer/
│       └── main.go                      # Entry point
├── internal/                            # Package imports
│   ├── server-pool/
│   │   ├── pool.go                      # Core logic, add/remove servers, track connections, health status
│   │   └── ...                          # Additional utility
│   └── lb-algorithms/                   # Implementations of load balancing algorithms
│       ├── round_robin.go
│       ├── weighted_round_robin.go
│       ├── least_connections.go
│       ├── ip_hash.go
│       ├── random.go
│       ├── ...                          # Any other algorithms
│       └── utils.go                     # Common helpers
├── demo/                                # Simulate usage of the LB
│   ├── client.go                        # Send client packets
│   └── server.go                        # Simple server template
├── docker/                              # Contains all files retaining to the docker image and containers
│   ├── Dockerfile.loadbalancer          # Builds binary image for main.go
│   └── Dockerfile.server                # Builds binary image for server.go                      
├── servers.json                         # Server list (file type tbd)                           
├── .dockerignore                        # Similar to .gitignore, outlines files for the docker container to ignore
├── compose.yml                          # File used for building the docker containers and describing how they should interact with each other
├── go.mod                               # Go module definition (tracks dependencies and package name) 
└── README.MD                            # ...
```
