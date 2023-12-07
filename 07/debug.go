package main

import "strings"

func scoreToString(score int) string {
	res := strings.Builder{}
	res.WriteString(handToString(hand(score >> 20)))
	res.WriteString(" ")
	res.WriteString(cardToString(card(score >> 16 & 0xF)))
	res.WriteString(cardToString(card(score >> 12 & 0xF)))
	res.WriteString(cardToString(card(score >> 8 & 0xF)))
	res.WriteString(cardToString(card(score >> 4 & 0xF)))
	res.WriteString(cardToString(card(score & 0xF)))
	return res.String()
}

func handToString(h hand) string {
	switch h {
	case fiveOfAKind:
		return "Five of a kind"
	case fourOfAKind:
		return "Four of a kind"
	case fullHouse:
		return "Full house"
	case threeOfAKind:
		return "Three of a kind"
	case twoPair:
		return "Two pair"
	case onePair:
		return "One pair"
	default:
		return "High card"
	}
}

func cardToString(c card) string {
	switch c {
	case ace:
		return "A"
	case king:
		return "K"
	case queen:
		return "Q"
	case jack:
		return "J"
	case joker:
		return "J"
	case ten:
		return "T"
	default:
		return string('1' + c)
	}
}
