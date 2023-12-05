package main

import (
	"github.com/csmith/aoc-2023/common"
	"github.com/csmith/aoc-2023/common/channels"
	"math"
	"strconv"
	"strings"
)

func main() {
	channels.RunForked(
		buildMaps(convertNumbers(parseLines(channels.SplitLines(channels.ReadFileBytes("05/input.txt"))))),
		partOne,
		partTwo,
	)
}

func partOne(in <-chan common.Pair[[]int, []*mapping]) <-chan int {
	return channels.Min(singleSeedLocations(in))
}

func partTwo(in <-chan common.Pair[[]int, []*mapping]) <-chan int {
	return channels.Min(rangeSeedLocations(in))
}

func singleSeedLocations(in <-chan common.Pair[[]int, []*mapping]) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		pair := <-in
		seeds := pair.First
		mappings := pair.Second

		for _, s := range seeds {
			for _, m := range mappings {
				if s >= m.from && s < m.from+m.length {
					out <- s - m.from + m.to
					break
				}
			}
		}
	}()

	return out
}

func rangeSeedLocations(in <-chan common.Pair[[]int, []*mapping]) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		pair := <-in
		seeds := pair.First
		mappings := pair.Second

		fakeSeedMappings := make([]*mapping, len(seeds)/2)
		for i := 0; i < len(seeds); i += 2 {
			fakeSeedMappings[i/2] = &mapping{
				from:   seeds[i],
				to:     seeds[i],
				length: seeds[i+1],
			}
		}

		for _, s := range fakeSeedMappings {
			for _, m := range mappings {
				o, _ := s.merge(m)
				if o != nil {
					out <- o.to
				}
			}
		}
	}()

	return out
}

func buildMaps(in <-chan []int) <-chan common.Pair[[]int, []*mapping] {
	out := make(chan common.Pair[[]int, []*mapping])
	go func() {
		defer close(out)
		seeds := <-in
		var mergedMaps = []*mapping{
			{
				from:   0,
				to:     0,
				length: math.MaxInt64,
			},
		}
		var currentMaps []*mapping

		for line := range in {
			if len(line) == 3 {
				// An actual mapping!
				currentMaps = append(currentMaps, &mapping{
					to:     line[0],
					from:   line[1],
					length: line[2],
				})
			} else if len(currentMaps) > 0 {
				if len(mergedMaps) == 0 {
					mergedMaps = currentMaps
				} else {
					// Time to merge!
					var newMaps []*mapping
					for _, m := range currentMaps {
						for i := range mergedMaps {
							if mergedMaps[i].consumed {
								continue
							}

							newMap, remaining := mergedMaps[i].merge(m)
							if newMap != nil {
								newMaps = append(newMaps, newMap)
								mergedMaps = append(mergedMaps, remaining...)
								mergedMaps[i].consumed = true
							}
						}
					}

					for _, m := range mergedMaps {
						if !m.consumed && m.length > 0 {
							newMaps = append(newMaps, m)
						}
					}

					mergedMaps = newMaps
				}

				currentMaps = currentMaps[:0]
			}
		}

		out <- common.Pair[[]int, []*mapping]{First: seeds, Second: mergedMaps}
	}()

	return out
}

func convertNumbers(in <-chan []string) <-chan []int {
	return channels.Map(in, func(line []string) []int {
		res := make([]int, len(line))
		for i, v := range line {
			res[i], _ = strconv.Atoi(v)
		}
		return res
	})
}

func parseLines(in <-chan string) <-chan []string {
	out := make(chan []string)
	go func() {
		defer close(out)

		// First line is a special snowflake
		out <- strings.Split(strings.TrimPrefix(<-in, "seeds: "), " ")

		for line := range in {
			out <- strings.Split(line, " ")
		}
	}()
	return out
}

type mapping struct {
	from     int
	to       int
	length   int
	consumed bool
}

func (m *mapping) merge(other *mapping) (new *mapping, remaining []*mapping) {
	otherStart := other.from
	otherEnd := other.from + other.length
	myStart := m.to
	myEnd := m.to + m.length

	if otherStart > myStart && otherStart <= myEnd {
		// There's a bit of our mapping at the start that remains unchanged
		remaining = append(remaining, &mapping{
			from:   m.from,
			to:     m.to,
			length: otherStart - myStart,
		})
	}

	if otherEnd >= myStart && otherEnd < myEnd {
		// There's a bit of our mapping at the end that remains unchanged
		overlap := otherEnd - myStart
		remaining = append(remaining, &mapping{
			from:   m.from + overlap,
			to:     m.to + overlap,
			length: m.length - overlap,
		})
	}

	if otherEnd >= myStart && otherStart <= myEnd {
		// There's some overlap that can become a shiny new merged mapping
		fromOffset := max(0, otherStart-myStart)
		toOffset := max(0, myStart-otherStart)
		length := min(other.length-toOffset, m.length-fromOffset)

		new = &mapping{
			from:   m.from + fromOffset,
			to:     other.to + toOffset,
			length: length,
		}
	}

	return
}
