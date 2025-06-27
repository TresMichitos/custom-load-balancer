[![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev/)

# Custom Load Balancer Using Go

Custom implementation of a load balancer using Go.

## Installation

1. Download the Repository with:

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
  go run .
  ```

## Development

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

### Possible Libraries To Use

- Included in Go standard library:
  - net/http
  - sync

