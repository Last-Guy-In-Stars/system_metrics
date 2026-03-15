package main

import (
	"context"
	"log"
	"project/project/proto"
	"sync"
	"time"

	"google.golang.org/grpc"
)

func main() {
	addresses := []string{
		// "192.168.1.97:50051",
		// "192.168.1.69:50051",
		"192.168.1.68:50051", // add here your ip addresses
	}
	timeout := 2 * time.Second
	for {
		var wg sync.WaitGroup
		for i := range addresses {
			wg.Add(1)

			go func(addr string) {
				defer wg.Done()
				metrics, err := getMetrics(addr, timeout)
				if err != nil {
					log.Println("Error:", err)
				} else {
					log.Printf(
						"CPU: %d%%, Memory: %dMB, Hostname: %s, OS: %s, Temperature: %dC\n",
						metrics.CpuUsage,
						metrics.MemoryUsage,
						metrics.OsName,
						metrics.Platform,
						metrics.Temperature)
				}
			}(addresses[i])
		}
		wg.Wait()
		time.Sleep(1 * time.Second)
	}
}

func getMetrics(address string, timeout time.Duration) (*proto.MetricsResponse, error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := proto.NewAgentServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return client.GetMetrics(ctx, &proto.EmptyRequest{})
}
