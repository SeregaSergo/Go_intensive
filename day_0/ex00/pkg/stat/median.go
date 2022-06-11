package stat

import (
	"math"
	"sort"
)

func Median(inputArray []float64) float64 {
	var length int
	length = len(inputArray)
	if length == 0 {
		return math.NaN()
	}
	sort.Float64s(inputArray)
	if length%2 == 0 {
		return float64(inputArray[length/2]+inputArray[(length-1)/2]) / 2
	} else {
		return inputArray[length/2]
	}
}
