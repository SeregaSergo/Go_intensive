package stat

import (
	"math"
)

func Deviation(inputArray []float64) float64 {
	var result float64
	mean := Mean(inputArray)
	for _, num := range inputArray {
		result = result + math.Pow(num-mean, 2)
	}
	return math.Sqrt(result / float64(len(inputArray)))
}
