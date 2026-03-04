package main

import (
	"fmt"
	"log"
	"net"
	"google.golang.org/grpc"
	"project/project/proto"
)

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
    for _, address := range addrs {
        // check the address type and if it is not a loopback the display it
        if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                return ipnet.IP.String()
            }
        }
    }
	return "0.0.0.0"
}

func main() {
	LocalIP := GetLocalIP()
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		fmt.Println("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAgentServiceServer(grpcServer, &server{})

	log.Printf("Agent listening at %s:50051", LocalIP)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}