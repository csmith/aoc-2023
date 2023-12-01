package common

// Last returns the last value received from the channel before it was closed
func Last(channel <-chan int) (res int) {
	for {
		o, more := <-channel
		if more {
			res = o
		} else {
			return
		}
	}
}

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

// Fork reads all data from the given channel and sends it to two new channels.
func Fork[T any](channel <-chan T) (chan T, chan T) {
	a := make(chan T, 100)
	b := make(chan T, 100)
	go func() {
		defer close(a)
		defer close(b)
		for o := range channel {
			a <- o
			b <- o
		}
	}()
	return a, b
}

// RunForkedChannels takes an input channel, forks it, and passes the values to the given functions.
// The first value from each function is printed to stdout.
func RunForkedChannels[T any, R any](input <-chan T, partOne func(<-chan T) <-chan R, partTwo func(<-chan T) <-chan R) {
	inputA, inputB := Fork(input)
	outputA := partOne(inputA)
	outputB := partTwo(inputB)
	println(<-outputA)
	println(<-outputB)
}
