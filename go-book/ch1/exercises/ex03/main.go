package main

import (
	"fmt"
	"math"
	"strings"
	"time"
)

func main() {
	const LettersInEnglish = 26
	const MaxArguments = 10000
	const MaxTestRuns = 10000
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	var arguments [MaxArguments]string
	var concatTimes [MaxTestRuns]float64
	var concatTimesSum, concatTimesMax, concatTimesMin float64
	var joinTimes [MaxTestRuns]float64
	var joinTimesSum, joinTimesMax, joinTimesMin float64

	prefix := ""
	prefixIndex := -1
	for iteration := 0; iteration < MaxArguments; iteration++ {
		for i := 0; i < LettersInEnglish; i++ {
			arguments[iteration] = prefix + string(letters[i])
			iteration++
			if iteration >= MaxArguments {
				break
			}
		}

		iteration--
		prefixIndex++
		prefix = arguments[prefixIndex]
	}

	// Execute Test runs
	var s, sep string

	concatTimesSum = 0
	concatTimesMax = -math.MaxFloat64
	concatTimesMin = math.MaxFloat64

	joinTimesSum = 0
	joinTimesMax = -math.MaxFloat64
	joinTimesMin = math.MaxFloat64

	for t := 0; t < MaxTestRuns; t++ {
		// Concatenation times
		start := time.Now()
		for i := 1; i < MaxArguments; i++ {
			s += sep + arguments[i]
			sep = " "
		}

		lap := time.Now()

		// get concat values
		elapsed := lap.Sub(start)
		var c float64 = float64(float64(elapsed.Microseconds()) / float64(1000))
		concatTimes[t] = c
		concatTimesSum += c
		if c < concatTimesMin {
			concatTimesMin = c
		}

		if c > concatTimesMax {
			concatTimesMax = c
		}

		// Join times
		start = time.Now()
		s = strings.Join(arguments[:], " ")
		lap = time.Now()

		// get join values
		elapsed = lap.Sub(start)
		var j float64 = float64(float64(elapsed.Microseconds()) / float64(1000))
		joinTimes[t] = j
		joinTimesSum += j
		if j < joinTimesMin {
			joinTimesMin = j
		}

		if j > joinTimesMax {
			joinTimesMax = j
		}

		fmt.Printf("%v", t%10)
		if t%100 == 0 {
			fmt.Println()
			fmt.Println(t)
			fmt.Printf("Avg Concat time: %v millisecs\n", float64(concatTimesSum/float64(t)))
			fmt.Printf("Avg Join   time: %v millisecs\n", float64(joinTimesSum/float64(t)))
			fmt.Println()

			fmt.Printf("Max Concat time: %v millisecs\n", float64(concatTimesMax))
			fmt.Printf("Max Join   time: %v millisecs\n", float64(joinTimesMax))
			fmt.Println()

			fmt.Printf("Min Concat time: %v millisecs\n", float64(concatTimesMin))
			fmt.Printf("Min Join   time: %v millisecs\n", float64(joinTimesMin))
			fmt.Println()
		}
	}

	fmt.Println()

	fmt.Printf("Avg Concat time: %v millisecs\n", float64(concatTimesSum/MaxTestRuns))
	fmt.Printf("Avg Join   time: %v millisecs\n", float64(joinTimesSum/MaxTestRuns))
	fmt.Println()

	fmt.Printf("Max Concat time: %v millisecs\n", float64(concatTimesMax))
	fmt.Printf("Max Join   time: %v millisecs\n", float64(joinTimesMax))
	fmt.Println()

	fmt.Printf("Min Concat time: %v millisecs\n", float64(concatTimesMin))
	fmt.Printf("Min Join   time: %v millisecs\n", float64(joinTimesMin))
	fmt.Println()
}
