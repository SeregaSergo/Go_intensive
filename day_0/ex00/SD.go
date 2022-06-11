package main

import "math"

func Deviation(inputArray []int) float64 {
	var result float64
	mean := Mean(inputArray)
	for _, num := range inputArray {
		result = result + math.Pow(float64(num)-mean, 2)
	}
	return math.Sqrt(result / float64(len(inputArray)))
}
