package channels

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
