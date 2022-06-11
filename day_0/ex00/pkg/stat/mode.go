package stat

import "math"

func Mode(inputArray []float64) float64 {
	if len(inputArray) == 0 {
		return math.NaN()
	}
	m := make(map[float64]int)
	for _, num := range inputArray {
		m[num] = m[num] + 1
	}
	var result float64
	var times int
	for k, v := range m {
		if v < times || (v == times && k > result) {
			continue
		} else {
			result = k
			times = v
		}
	}
	return result
}
