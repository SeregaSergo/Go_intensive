package main

import (
	"bufio"
	"ex00/pkg/stat"
	"flag"
	"fmt"
	"os"
	"strconv"
)

func main() {
	flags := getFlags()
	inputArray := make([]float64, 0)
	getInput(&inputArray)

	if *flags["mean"] {
		fmt.Printf("Mean: %v\n", stat.Mean(inputArray))
	}

	if *flags["median"] {
		fmt.Printf("Median: %v\n", stat.Median(inputArray))
	}

	if *flags["mode"] {
		fmt.Printf("Mode: %v\n", stat.Mode(inputArray))
	}

	if *flags["deviation"] {
		fmt.Printf("SD: %v\n", stat.Deviation(inputArray))
	}
}

func getFlags() map[string]*bool {
	m := make(map[string]*bool)
	m["mode"] = flag.Bool("mode", false, "count mode")
	m["mean"] = flag.Bool("mean", false, "count mean")
	m["median"] = flag.Bool("median", false, "count median")
	m["deviation"] = flag.Bool("deviation", false, "count standard deviation")
	flag.Parse()

	if checkIfNoFlags(m) {
		changeFlagsToTrue(m)
	}
	return m
}

func checkIfNoFlags(m map[string]*bool) bool {
	numParam := 0
	for _, v := range m {
		if *v == false {
			numParam = numParam + 1
		}
	}
	return numParam == 4
}

func changeFlagsToTrue(m map[string]*bool) {
	for _, v := range m {
		*v = true
	}
}

func getInput(inputArray *[]float64) {
	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		line := sc.Text()
		if s, err := strconv.Atoi(line); err != nil {
			e := err.(*strconv.NumError).Err
			switch {
			case e != strconv.ErrSyntax:
				fmt.Println("Error:", e)
			case len(line) == 0:
				fmt.Println("Error: empty line")
			default:
				fmt.Println("Error: is not a number")
			}
		} else {
			*inputArray = append(*inputArray, float64(s))
		}
	}
}
