package channels

import (
	"errors"
	"io"
	"os"
)

// ReadFileBytes reads all bytes from the given path and emits them in a
// channel. A LF is always added to the end of the file to ease processing.
func ReadFileBytes(path string) <-chan byte {
	ch := make(chan byte)
	go func() {
		defer close(ch)
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer func() {
			_ = file.Close()
		}()

		b := make([]byte, 1024)
		for {
			n, err := file.Read(b)
			if errors.Is(err, io.EOF) {
				break
			}
			if err != nil {
				panic(err)
			}
			for i := 0; i < n; i++ {
				ch <- b[i]
			}
		}

		ch <- '\n'
	}()
	return ch
}
