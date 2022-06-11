package stat

func Mean(inputArray []float64) float64 {
	var result float64
	for _, num := range inputArray {
		result = result + num
	}
	return result / float64(len(inputArray))
}
