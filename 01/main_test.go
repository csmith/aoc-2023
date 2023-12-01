package main

import (
	"github.com/csmith/aoc-2023/common"
	"testing"
)

func BenchmarkDayOne(b *testing.B) {
	for i := 0; i < b.N; i++ {
		common.RunForkedChannels(common.ReadFileChannel("input.txt"), partOne, partTwo)
	}
}
