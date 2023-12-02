package main

import (
	"fmt"
	"github.com/csmith/aoc-2023/common"
	"github.com/csmith/aoc-2023/common/channels"
	"strings"
)

func main() {
	channels.RunForked(parseInput(), partOne, partTwo)
}

func partOne(in <-chan game) <-chan int {
	return sumGameIds(filterGames(in, 12, 13, 14))
}

func partTwo(in <-chan game) <-chan int {
	return channels.Sum(channels.Map(in, common.Compose2(minimumCubes, calculatePower)))
}

func parseInput() chan game {
	return channels.Map(
		channels.SplitLines(channels.ReadFileBytes("02/input.txt")),
		common.Compose6(
			parseGameId,
			splitRounds,
			splitCounts,
			parseCounts,
			aggregateCounts,
			createGame,
		),
	)
}

func parseGameId(line string) common.Pair[int, string] {
	var id int
	header, tail, _ := strings.Cut(line, ": ")
	_, _ = fmt.Sscanf(header, "Game %d", &id)
	return common.Pair[int, string]{First: id, Second: tail}
}

func splitRounds(in common.Pair[int, string]) common.Pair[int, []string] {
	return common.MapSecond(in, func(s string) []string {
		return strings.Split(s, "; ")
	})
}

func splitCounts(in common.Pair[int, []string]) common.Pair[int, [][]string] {
	return common.MapSecondSlice(in, func(s string) []string {
		return strings.Split(s, ", ")
	})
}

func parseCounts(in common.Pair[int, [][]string]) common.Pair[int, [][]selection] {
	return common.MapSecondSliceSlice(in, func(s string) selection {
		var count int
		var colour string
		var res selection

		_, _ = fmt.Sscanf(s, "%d %s", &count, &colour)

		switch colour {
		case "red":
			res.red = count
		case "green":
			res.green = count
		case "blue":
			res.blue = count
		}

		return res
	})
}

func aggregateCounts(in common.Pair[int, [][]selection]) common.Pair[int, []selection] {
	return common.MapSecondSlice(in, func(s []selection) selection {
		var res selection
		for _, sel := range s {
			res.red += sel.red
			res.green += sel.green
			res.blue += sel.blue
		}
		return res
	})
}

func createGame(in common.Pair[int, []selection]) game {
	return game{id: in.First, rounds: in.Second}
}

func filterGames(in <-chan game, maxRed, maxGreen, maxBlue int) <-chan game {
	return channels.Filter(in, func(g game) bool {
		for _, r := range g.rounds {
			if r.red > maxRed || r.green > maxGreen || r.blue > maxBlue {
				return false
			}
		}
		return true
	})
}

func sumGameIds(in <-chan game) <-chan int {
	return channels.Sum(channels.Map(in, func(g game) int {
		return g.id
	}))
}

func minimumCubes(g game) selection {
	var res selection
	for _, r := range g.rounds {
		res.red = max(res.red, r.red)
		res.green = max(res.green, r.green)
		res.blue = max(res.blue, r.blue)
	}
	return res
}

func calculatePower(s selection) int {
	return s.red * s.green * s.blue
}

type game struct {
	id     int
	rounds []selection
}

type selection struct {
	red   int
	green int
	blue  int
}
