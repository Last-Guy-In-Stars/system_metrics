package main

import (
	"context"
	"log"
	"time"

	"project/project/proto"

	"google.golang.org/grpc"
)

func main() {
	addressess := []string{
		"192.168.1.97:50051",
		"192.168.1.69:50051",
	}
	timeout := 12 * time.Second

	for {
		for i := range addressess {
			metrics, err := getMetrics(addressess[i], timeout)
			if err != nil {
				log.Println("Error:", err)
			} else {
				log.Printf("CPU: %d%%, Memory: %dMB, Hostname: %s, OS: %s\n",
					metrics.CpuUsage, metrics.MemoryUsage, metrics.OsName, metrics.Platform)
			}
			time.Sleep(4 * time.Second)
		}
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
