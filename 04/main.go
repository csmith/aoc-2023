package main

import (
	"github.com/csmith/aoc-2023/common/channels"
	"strings"
)

func main() {
	channels.RunForked(
		countMatches(parseCards(channels.SplitLines(channels.ReadFileBytes("04/input.txt")))),
		partOne,
		partTwo,
	)
}

func partOne(in <-chan int) <-chan int {
	return channels.Sum(calculatePoints(in))
}

func partTwo(in <-chan int) <-chan int {
	return channels.Sum(countCopies(in))
}

func countCopies(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)

		copies := 1
		var future []int

		for i := range in {
			out <- copies

			if len(future) < i {
				future = append(future, make([]int, i)...)
			}

			for j := 0; j < i; j++ {
				future[j] += copies
			}

			if len(future) > 0 {
				copies = 1 + future[0]
				future = future[1:]
			} else {
				copies = 1
			}
		}
	}()
	return out
}

func calculatePoints(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for i := range in {
			if i > 0 {
				out <- 1 << (i - 1)
			}
		}
	}()
	return out
}

func countMatches(in <-chan card) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for c := range in {
			winning := make(map[string]bool)
			for _, n := range c.winningNumbers {
				winning[n] = true
			}

			matches := 0
			for _, n := range c.cardNumbers {
				// Due to the weird formatting we have some "blank" numbers,
				// just ignore them here.
				if len(n) > 0 && winning[n] {
					matches++
				}
			}

			out <- matches
		}
	}()
	return out
}

func parseCards(in <-chan string) <-chan card {
	out := make(chan card)
	go func() {
		defer close(out)
		for line := range in {
			_, numbers, _ := strings.Cut(line, ": ")
			if len(numbers) > 0 {
				ours, winning, _ := strings.Cut(numbers, " | ")
				out <- card{
					cardNumbers:    strings.Split(ours, " "),
					winningNumbers: strings.Split(winning, " "),
				}
			}
		}
	}()
	return out
}

type card struct {
	cardNumbers    []string
	winningNumbers []string
}
