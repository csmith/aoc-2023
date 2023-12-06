package channels

import "fmt"

func Print[T any](channel <-chan T) <-chan T {
	out := make(chan T)
	go func() {
		defer close(out)
		i := 0
		for o := range channel {
			i++
			fmt.Printf("Item %d: %#v\n", i, o)
			out <- o
		}
	}()
	return out
}
