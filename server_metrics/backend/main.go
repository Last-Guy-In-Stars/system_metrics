package main

import (
	"context"
	"log"
	"sync"
	"time"

	"server_metrics/project/proto"

	"google.golang.org/grpc"
)

func main() {
	addresses := []string{
		"192.168.1.97:50051",
		"192.168.1.69:50051", // add here your ip addresses
	}
	timeout := 2 * time.Second
	for {
		var wg sync.WaitGroup // Create a WaitGroup to wait for all goroutines to finish
		for i := range addresses {
			wg.Add(1) //  Add 1 to the WaitGroup for each goroutine

			go func(addr string) { // Create a goroutine to fetch metrics from the server
				defer wg.Done() // Decrement the WaitGroup when the goroutine finishes
				metrics, err := getMetrics(addr, timeout)
				if err != nil {
					log.Println("Error:", err)
				} else {
					log.Printf(
						"CPU: %d%%, Memory: %dMB, Hostname: %s, OS: %s\n",
						metrics.CpuUsage,
						metrics.MemoryUsage,
						metrics.OsName,
						metrics.Platform)
				}
			}(addresses[i])
		}
		wg.Wait() // Wait for all goroutines to finish
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
