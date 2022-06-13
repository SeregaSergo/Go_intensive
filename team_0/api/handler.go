package api

import (
	"github.com/google/uuid"
	"gonum.org/v1/gonum/stat/distuv"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"math"
	"math/rand"
	"time"
)

type Transmitter struct{}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func (*Transmitter) mustEmbedUnimplementedCommunicationServer() {}

func (t *Transmitter) StreamFrequencies(in *Empty, stream Communication_StreamFrequenciesServer) error {
	rand.Seed(time.Now().Unix())
	id := uuid.New().String()
	dist := distuv.Normal{
		Mu:    roundFloat(float64(rand.Intn(20)-10)+rand.Float64(), 2),
		Sigma: rand.Float64()*1.2 + 0.3,
	}
	log.Printf("New connection: %s\nMean: %f\nStandard deviation: %f\n\n", id, dist.Mu, dist.Sigma)
	for {
		frequency := Frequency{
			SessionId: id,
			Frequency: dist.Rand(),
			Time:      timestamppb.Now(),
		}
		if err := stream.Send(&frequency); err != nil {
			return err
		}
		//time.Sleep(time.Second)
	}
}
