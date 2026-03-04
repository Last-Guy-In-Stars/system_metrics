package main

import (
	"log"
	"context"
	"google.golang.org/grpc/peer"
	"project/project/proto"
)

type server struct {
	proto.UnimplementedAgentServiceServer // Заглушка
}

func (s *server) GetMetrics(ctx context.Context, req *proto.EmptyRequest) (*proto.MetricsResponse, error) {
	// Логируем подключение сервера
	if p, ok := peer.FromContext(ctx); ok {
		log.Printf("Server connected: %v", p.Addr)
	} else {
		log.Println("Server connected: unknown")
	}
	// Получаем метрики
	return &proto.MetricsResponse{
		CpuUsage: int32(GetCPU()),
		MemoryUsage: int32(GetMemory()),
		OsName: string(GetOsName()),
	}, nil
}