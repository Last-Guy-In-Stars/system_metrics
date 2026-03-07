package main

import (
	"context"
	"log"
	"time"

	"server_metrics/project/proto"

	"google.golang.org/grpc"
)

func main() {
	addresses := []string{
		"IP:50051", // add here your ip addresses
	}
	timeout := 12 * time.Second

	for {
		for i := range addresses {
			metrics, err := getMetrics(addresses[i], timeout)
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
