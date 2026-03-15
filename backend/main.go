package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"project/project/proto"
	"sync"
	"time"

	"github.com/go-gota/gota/dataframe"
	"github.com/go-gota/gota/series"
	"google.golang.org/grpc"
)

var (
	addresses = []string{
		"192.168.1.97:50051",
		"192.168.1.69:50051",
		"192.168.1.68:50051", // add here your ip addresses
	}
	timeout = 10 * time.Second
)

func main() {
	go scheduleEmail()
	for {
		CollectAndSaveMetrics()
		time.Sleep(10 * time.Second)
	}
}

func CollectAndSaveMetrics() {
	var wg sync.WaitGroup
	results := make([]*proto.MetricsResponse, len(addresses))

	for i := range addresses {
		wg.Add(1)

		go func(addr string) {
			defer wg.Done()
			metrics, err := getMetrics(addr, timeout)
			if err != nil {
				log.Println("Error:", err)
			} else {
				results[i] = metrics
			}
		}(addresses[i])
	}
	wg.Wait()

	var names []string
	var cpu []int
	var mem []int
	var os_name []string
	var temp []int

	for _, m := range results {
		if m != nil {
			names = append(names, m.OsName)
			cpu = append(cpu, int(m.CpuUsage))
			mem = append(mem, int(m.MemoryUsage))
			os_name = append(os_name, m.Platform)
			temp = append(temp, int(m.Temperature))
		}
	}

	df := dataframe.New(
		series.New(names, series.String, "NamePC"),
		series.New(cpu, series.Int, "CPULoad"),
		series.New(mem, series.Int, "MemLoad"),
		series.New(os_name, series.String, "OS"),
		series.New(temp, series.Int, "Temp"),
	)

	filename := fmt.Sprintf("metrics_%s.csv", time.Now().Format("2007-07-07"))
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	df.WriteCSV(file)

	log.Println(df)

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
