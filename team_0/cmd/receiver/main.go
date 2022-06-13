package main

import (
	"context"
	"fmt"
	math_stat "github.com/SeregaSergo/Go_intensive/day_0/ex00/pkg/stat"
	pb "github.com/SeregaSergo/Go_intensive/team_0/api"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"io"
	"log"
)

func main() {

	var conn *grpc.ClientConn
	conn, err := grpc.Dial("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	defer conn.Close()

	client := pb.NewCommunicationClient(conn)

	stream, err := client.StreamFrequencies(context.Background(), &pb.Empty{})
	if err != nil {
		log.Fatalf("Error when calling StreamFrequencies: %s", err)
	}

	var frequencies []float64
	var mean, SD float64
	var i int
	for i := 0; {
		frequency, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("StreamFrequencies() =, %v", err)
		}
		i++
		frequencies = append(frequencies, frequency.Frequency)
		mean = math_stat.Mean(frequencies)
		SD = math_stat.Deviation(frequencies)
		//fmt.Println(time.Unix(frequency.GetTime().Seconds, 0).Format("[2006-01-02 15:04:05]"), frequency.GetSessionId()[:8])
		fmt.Println(i, "|", frequency.GetFrequency())
		fmt.Println("Approximated mean:", mean)
		fmt.Println("Approximated SD:", SD)
	}
}
