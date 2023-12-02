package channels

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

// RunForked takes an input channel, forks it, and passes the values to the given functions.
// The first value from each function is printed to stdout.
func RunForked[T any, R any](input <-chan T, partOne func(<-chan T) <-chan R, partTwo func(<-chan T) <-chan R) {
	inputA, inputB := Fork(input)
	outputA := partOne(inputA)
	outputB := partTwo(inputB)
	println(<-outputA)
	println(<-outputB)
}
