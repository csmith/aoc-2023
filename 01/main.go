package main

import (
	"github.com/csmith/aoc-2023/common/channels"
)

func main() {
	channels.RunForked(channels.ReadFileBytes("01/input.txt"), partOne, partTwo)
}

func partOne(in <-chan byte) <-chan int {
	return channels.Sum(produceSums(in))
}

func partTwo(in <-chan byte) <-chan int {
	return channels.Sum(produceSums(emitNumberWords(in)))
}

func produceSums(in <-chan byte) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)

		first := byte('?')
		last := byte('?')
		for b := range in {
			if b >= '0' && b <= '9' {
				last = b
				if first == '?' {
					first = b
				}
			}

			if b == '\n' && first != '?' {
				out <- int(last-'0') + int(first-'0')*10
				first = '?'
			}
		}
	}()
	return out
}

var numbers = []uint64{
	0x0000FFFFFF, 0x6F6E65,
	0x0000FFFFFF, 0x74776F,
	0xFFFFFFFFFF, 0x7468726565,
	0x00FFFFFFFF, 0x666F7572,
	0x00FFFFFFFF, 0x66697665,
	0x0000FFFFFF, 0x736978,
	0xFFFFFFFFFF, 0x736576656E,
	0xFFFFFFFFFF, 0x6569676874,
	0x00FFFFFFFF, 0x6E696E65,
}

func emitNumberWords(in <-chan byte) <-chan byte {
	out := make(chan byte)
	go func() {
		defer close(out)

		var buffer uint64 = 0
		for b := range in {
			out <- b
			buffer = (buffer << 8) | uint64(b)

			for i := 0; i < len(numbers); i += 2 {
				if buffer&numbers[i] == numbers[i+1] {
					out <- byte('1' + i/2)
				}
			}
		}
	}()
	return out
}
