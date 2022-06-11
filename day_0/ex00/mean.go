package main

func Mean(inputArray []int) float64 {
	var result float64
	for _, num := range inputArray {
		result = result + float64(num)
	}
	return result / float64(len(inputArray))
}
