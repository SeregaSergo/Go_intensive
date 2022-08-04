package main

import (
	"bufio"
	"context"
	_ "database/sql"
	"fmt"
	mathStat "github.com/SeregaSergo/Go_intensive/day_0/ex00/pkg/stat"
	"github.com/SeregaSergo/Go_intensive/team_0/internal/storage"
	"github.com/SeregaSergo/Go_intensive/team_0/internal/transport"
	"github.com/dariubs/percent"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"math"
	"os"
	"strconv"
)

const (
	DB_HOST     = "localhost"
	DB_PORT     = 5432
	DB_USER     = "postgres"
	DB_PASSWORD = "kek"
	DB_NAME     = "serega"

	approximationRate = 0.01
	address           = ":9000"
)

func main() {
	conn := getConnection()
	defer conn.Close()

	stream := getStream(address, conn)
	mean, SD := approximateStats(stream)
	k := getAnomalyRate()
	min, max := getAnomalyLimits(mean, SD, k)
	db := connectDB()

	frequencyMsg, err := stream.Recv()
	for err == nil {
		if frequencyMsg.Frequency < min || frequencyMsg.Frequency > max {
			result := db.Create(&storage.FrequencyRecord{
				SessionId: frequencyMsg.SessionId,
				Frequency: frequencyMsg.Frequency,
				Timestamp: fmt.Sprintf("%d", frequencyMsg.Time.Seconds),
			})
			fmt.Println(result.RowsAffected, " Row is recorded")
		}
		frequencyMsg, err = stream.Recv()
	}
	var record storage.FrequencyRecord
	db.First(&record)
	fmt.Println(record)
}

func getConnection() *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	return conn
}

func getStream(address string, conn *grpc.ClientConn) transport.Communication_StreamFrequenciesClient {
	client := transport.NewCommunicationClient(conn)
	stream, err := client.StreamFrequencies(context.Background(), &transport.Empty{})
	if err != nil {
		log.Fatalf("Error when calling StreamFrequencies: %s", err)
	}
	return stream
}

func approximateStats(stream transport.Communication_StreamFrequenciesClient) (mean float64, SD float64) {
	var frequencies []float64
	var counter int
	var tempMean, tempSD float64
	mean = 1.0
	SD = 1.0
	for math.Abs(percent.ChangeFloat(tempMean, mean)) > approximationRate ||
		math.Abs(percent.ChangeFloat(tempSD, SD)) > approximationRate {
		tempMean = mean
		tempSD = SD

		frequency, err := stream.Recv()
		if err != nil {
			log.Fatalf("StreamFrequencies() =, %v", err)
		}
		frequencies = append(frequencies, frequency.Frequency)

		counter++
		mean = mathStat.Mean(frequencies)
		SD = mathStat.Deviation(frequencies)
		fmt.Println(counter, "frequencies were received (", frequency, ")")
		fmt.Println("Approximated mean:", mean)
		fmt.Println("Approximated SD:", SD)
		fmt.Println("")
	}
	return
}

func getAnomalyRate() (k float64) {
	var err error
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter an anomaly coefficient (float value):")
	for scanner.Scan() {
		if k, err = strconv.ParseFloat(scanner.Text(), 64); err != nil {
			fmt.Println("ERROR: You need to enter one float number")
			continue
		}
		if k < 0 {
			fmt.Println("ERROR: Anomaly coefficient can't be negative")
			continue
		}
		break
	}
	return
}

func getAnomalyLimits(mean float64, SD float64, k float64) (min float64, max float64) {
	min = mean - SD*k
	max = mean + SD*k
	return
}

func connectDB() *gorm.DB {
	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		DB_HOST, DB_PORT, DB_USER, DB_PASSWORD, DB_NAME)
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Println("engine creation failed", err)
	}
	dbSql, _ := db.DB()
	err = dbSql.Ping()
	if err != nil {
		panic(err)
	}

	log.Println("Successfully connected")
	db.AutoMigrate(&storage.FrequencyRecord{})
	return db
}
