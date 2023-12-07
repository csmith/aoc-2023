package main

import (
	"github.com/csmith/aoc-2023/common"
	"github.com/csmith/aoc-2023/common/channels"
)

func main() {
	channels.RunForked(
		channels.ReadFileBytes("07/input.txt"),
		partOne,
		partTwo,
	)
}

func partOne(in <-chan byte) <-chan int {
	return channels.Sum(
		calculatePayouts(
			channels.Sort(
				parseHands(in, false),
				compareHands,
			),
		),
	)
}

func partTwo(in <-chan byte) <-chan int {
	return channels.Sum(
		calculatePayouts(
			channels.Sort(
				parseHands(in, true),
				compareHands,
			),
		),
	)
}

type card uint8

const (
	joker card = iota
	two
	three
	four
	five
	six
	seven
	eight
	nine
	ten
	jack
	queen
	king
	ace
	numCards
)

var cardSymbols = map[byte]card{
	'A': ace,
	'K': king,
	'Q': queen,
	'J': jack,
	'T': ten,
	'9': nine,
	'8': eight,
	'7': seven,
	'6': six,
	'5': five,
	'4': four,
	'3': three,
	'2': two,
}

type hand uint8

const (
	highCard hand = iota
	onePair
	twoPair
	threeOfAKind
	fullHouse
	fourOfAKind
	fiveOfAKind
)

func compareHands(a, b common.Pair[int, int]) int {
	return a.First - b.First
}

func calculatePayouts(in <-chan common.Pair[int, int]) <-chan int {
	out := make(chan int)

	go func() {
		defer close(out)

		rank := 1
		for p := range in {
			out <- p.Second * rank
			rank++
		}
	}()

	return out
}

func parseHands(in <-chan byte, jokers bool) <-chan common.Pair[int, int] {
	out := make(chan common.Pair[int, int])

	go func() {
		defer close(out)

		i := 1
		score := 0
		bid := 0
		var cardCounts [numCards]uint8
		var dupes [6]uint8
		for b := range in {
			if i <= 5 {
				c := cardSymbols[b]
				if c == jack && jokers {
					c = joker
				}
				score = score<<4 | int(c)
				cardCounts[c]++
				dupes[cardCounts[c]-1]--
				dupes[cardCounts[c]]++
			} else if b >= '0' && b <= '9' {
				bid = bid*10 + int(b-'0')
			} else if b == '\n' {
				jokerCount := cardCounts[joker]
				dupes[jokerCount]--
				score = int(jokerMutationMap[getHand(dupes)][jokerCount])<<20 | score
				out <- common.Pair[int, int]{First: score, Second: bid}
				score = 0
				bid = 0
				cardCounts = [numCards]uint8{}
				dupes = [6]uint8{}
				i = 0
			}

			i++
		}
	}()

	return out
}

const invalid = hand(255)

// jokerMutationMap maps each hand to a new hand based on the number of jokers present
var jokerMutationMap = [7][6]hand{
	// High card
	{highCard, onePair, threeOfAKind, fourOfAKind, fiveOfAKind, fiveOfAKind},
	// One pair
	{onePair, threeOfAKind, fourOfAKind, fiveOfAKind, invalid, invalid},
	// Two pair
	{twoPair, fullHouse, fourOfAKind, fiveOfAKind, invalid, invalid},
	// Three of a kind
	{threeOfAKind, fourOfAKind, fiveOfAKind, invalid, invalid, invalid},
	// Full house
	{fullHouse, fourOfAKind, fiveOfAKind, invalid, invalid, invalid},
	// Four of a kind
	{fourOfAKind, fiveOfAKind, invalid, invalid, invalid, invalid},
	// Five of a kind
	{fiveOfAKind, invalid, invalid, invalid, invalid, invalid},
}

func getHand(dupes [6]uint8) hand {
	switch {
	case dupes[5] == 1:
		return fiveOfAKind
	case dupes[4] == 1:
		return fourOfAKind
	case dupes[3] == 1 && dupes[2] == 1:
		return fullHouse
	case dupes[3] == 1:
		return threeOfAKind
	case dupes[2] == 2:
		return twoPair
	case dupes[2] == 1:
		return onePair
	default:
		return highCard
	}
}
