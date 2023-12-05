package channels

import "math"

// Sum adds up all the values received on the channel
func Sum(channel <-chan int) chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		sum := 0
		for o := range channel {
			sum += o
		}
		res <- sum
	}()
	return res
}

// Min finds the smallest value received on the channel
func Min(channel <-chan int) chan int {
	res := make(chan int)
	go func() {
		defer close(res)
		lowest := math.MaxInt64
		for o := range channel {
			if o < lowest {
				lowest = o
			}
		}
		res <- lowest
	}()
	return res
}
