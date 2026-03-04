package main

import (
	"context"
	"log"
	"time"

	"google.golang.org/grpc"
	"project/project/proto"
)

func main() {
	agentAddress := "192.168.1.69:50051" // IP:порт агента
	timeout := 12 * time.Second

	for {
		metrics, err := getMetrics(agentAddress, timeout)
		if err != nil {
			log.Println("Error:", err)
		} else {
			log.Printf("CPU: %d%%, Memory: %dMB, Hostname: %s\n", metrics.CpuUsage, metrics.MemoryUsage, metrics.OsName)
		}
		time.Sleep(8 * time.Second)
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