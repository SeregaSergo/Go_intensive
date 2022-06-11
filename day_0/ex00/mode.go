package main

import "math"

func Mode(inputArray []int) float64 {
	if len(inputArray) == 0 {
		return math.NaN()
	}
	m := make(map[int]int)
	for _, num := range inputArray {
		m[num] = m[num] + 1
	}
	var result int
	var times int
	for k, v := range m {
		if v < times || (v == times && k > result) {
			continue
		} else {
			result = k
			times = v
		}
	}
	return float64(result)
}
