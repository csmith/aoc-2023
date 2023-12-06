package main

import (
	"github.com/csmith/aoc-2023/common/channels"
	"math"
	"strconv"
	"strings"
)

func main() {
	channels.RunForked(
		channels.SplitLines(channels.ReadFileBytes("06/input.txt")),
		partOne,
		partTwo,
	)
}

func partOne(in <-chan string) <-chan int {
	return channels.Product(countOptions(zipLines(in)))
}

func partTwo(in <-chan string) <-chan int {
	return channels.Product(countOptions(zipLines(mergeObviouslySeparateNumbersIntoOne(in))))
}

func mergeObviouslySeparateNumbersIntoOne(in <-chan string) <-chan string {
	out := make(chan string)

	go func() {
		defer close(out)

		for line := range in {
			if len(line) > 12 {
				out <- line[0:12] + strings.ReplaceAll(line[12:], " ", "")
			}
		}
	}()

	return out
}

func zipLines(in <-chan string) <-chan [2]int {
	out := make(chan [2]int)

	go func() {
		defer close(out)

		first := strings.Fields(<-in)
		second := strings.Fields(<-in)

		for i := 1; i < len(first); i++ {
			f, _ := strconv.Atoi(first[i])
			s, _ := strconv.Atoi(second[i])
			out <- [2]int{f, s}
		}
	}()

	return out
}

func countOptions(in <-chan [2]int) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		for pair := range in {
			// distance = (time - button_press_time) * button_press_time
			// d = t * bpt - bpt^2
			// bpt^2 - t * bpt + d = 0
			// [Googles "quadratic equation"... It's been a long time.]
			// bpt = t +/- sqrt(t^2 - 4d) / 2

			// We want to _beat_ the distance, so solve for slightly longer.
			// This feels like a horrible hack that might break, and I don't
			// fully understand why I'm taking away the 0.000001 instead of
			// adding it, but it works…
			distance := float64(pair[0]) - 0.000001
			sqrt := math.Sqrt(distance*distance - 4*float64(pair[1]))
			lower := int(math.Ceil((distance - sqrt) / 2))
			upper := int(math.Floor((distance + sqrt) / 2))
			out <- 1 + upper - lower
		}
	}()

	return out
}
