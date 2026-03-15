package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"server_metrics/project/proto"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

func GetNewPort() string {
	const defaultPort = "50051"
	reader := bufio.NewReader(os.Stdin)

	fmt.Printf("Default port is %s\n", defaultPort)
	fmt.Print("Enter a new port number (1-65535) or press Enter to use default: ")

	for {
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text)

		if text == "" {
			return defaultPort
		}

		port, err := strconv.Atoi(text)
		if err != nil {
			fmt.Printf("Invalid port '%s'.\n", text)
			fmt.Print("Please enter a number: ")
			continue
		}

		if port < 1 || port > 65535 {
			fmt.Printf("Port %d is out of range (1-65535).\n",
				port)
			fmt.Print("Please enter a number: ")
			continue
		}

		fmt.Printf("Using port %d.\n", port)
		return ":" + text
	}
}

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
	newPort := GetNewPort()
	LocalIP := GetLocalIP()
	lis, err := net.Listen("tcp", newPort)
	if err != nil {
		log.Fatalf("failed to listen:", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterAgentServiceServer(grpcServer, &server{})

	log.Printf("Agent listening at %s%s", LocalIP, newPort)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
