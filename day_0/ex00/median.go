package main

import (
	"math"
	"sort"
)

func Median(inputArray []int) float64 {
	var length int
	length = len(inputArray)
	if length == 0 {
		return math.NaN()
	}
	sort.Ints(inputArray)
	if length%2 == 0 {
		return float64(inputArray[length/2]+inputArray[(length-1)/2]) / 2
	} else {
		return float64(inputArray[length/2])
	}
}