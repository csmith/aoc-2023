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

// SplitLines splits a byte channel into a string channel, splitting on LFs
// and ignoring CRs.
func SplitLines(in <-chan byte) <-chan string {
	out := make(chan string)
	go func() {
		defer close(out)
		var buffer []byte
		for b := range in {
			if b == '\n' {
				out <- string(buffer)
				buffer = nil
			} else if b != '\r' {
				buffer = append(buffer, b)
			}
		}
	}()
	return out
}
