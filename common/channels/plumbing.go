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
func RunForked[T, R any](input <-chan T, partOne func(<-chan T) <-chan R, partTwo func(<-chan T) <-chan R) {
	inputA, inputB := Fork(input)
	outputA := partOne(inputA)
	outputB := partTwo(inputB)
	println(<-outputA)
	println(<-outputB)
}

// Map applies the given function to each value received on the channel and
// emits it to the returned channel.
func Map[T, R any](channel <-chan T, f func(T) R) chan R {
	res := make(chan R)
	go func() {
		defer close(res)
		for o := range channel {
			res <- f(o)
		}
	}()
	return res
}

// Filter applies the given function to each value received on the channel and
// emits it to the returned channel if the function returns true.
func Filter[T any](channel <-chan T, f func(T) bool) chan T {
	res := make(chan T)
	go func() {
		defer close(res)
		for o := range channel {
			if f(o) {
				res <- o
			}
		}
	}()
	return res
}
