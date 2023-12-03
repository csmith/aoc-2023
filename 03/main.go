package main

import "github.com/csmith/aoc-2023/common/channels"

func main() {
	channels.RunForked(
		linkNeighbours(addBorders(parseSchematic(channels.ReadFileBytes("03/input.txt")))),
		partOne,
		partTwo,
	)
}

func partOne(in <-chan []cell) <-chan int {
	return channels.Sum(findPartNumbers(in))
}

func partTwo(in <-chan []cell) <-chan int {
	return channels.Sum(findGears(in))
}

func findPartNumbers(in <-chan []cell) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for row := range in {
			inNumber := false
			for _, c := range row {
				if c.is(cellKindNumber) {
					if inNumber {
						continue
					}
					inNumber = true
					cells, number := chaseNumber(&c)
					if len(findAdjacent(cells, cellKindSymbol)) >= 1 {
						out <- number
					}
				} else {
					inNumber = false
				}
			}
		}
	}()
	return out
}

func findGears(in <-chan []cell) <-chan int {
	out := make(chan int)
	go func() {
		defer close(out)
		for row := range in {
			for _, c := range row {
				if c.is(cellKindSymbol) && c.value == '*' {
					matches := findAdjacent([]*cell{&c}, cellKindNumber)
					if len(matches) == 2 {
						_, a := chaseNumber(matches[0])
						_, b := chaseNumber(matches[1])
						out <- a * b
					}
				}
			}
		}
	}()
	return out
}

// findStartOfNumber finds the first cell in a multi-cell number
func findStartOfNumber(c *cell) *cell {
	current := c
	for current.neighbours[left].is(cellKindNumber) {
		current = current.neighbours[left]
	}
	return current
}

// chaseNumber returns the cells and total value of a multi-cell number that
// contains the given cell
func chaseNumber(c *cell) ([]*cell, int) {
	start := findStartOfNumber(c)

	var cells []*cell
	var number = 0
	var current = start
	for current.is(cellKindNumber) {
		cells = append(cells, current)
		number = number*10 + current.value
		current = current.neighbours[right]
	}

	return cells, number
}

// findAdjacent finds all cells of the given kind that neighbour the given
// range. Cells must be horizontally adjacent and in order. Numbers are
// automatically de-duped, and only the starting position is returned.
func findAdjacent(cells []*cell, kind cellKind) []*cell {
	found := make(map[int]bool)
	var res []*cell

	for i, c := range cells {
		start := 0
		end := 8

		if i > 0 {
			// Don't look left for anything after the first cell
			start += 3
		}

		if i < len(cells)-1 {
			// Don't look right for anything but the last cell
			end -= 3
		}

		for _, n := range c.neighbours[start:end] {
			if n.is(kind) {
				resolved := n
				if resolved.is(cellKindNumber) {
					resolved = findStartOfNumber(resolved)
				}

				if !found[resolved.id] {
					found[resolved.id] = true
					res = append(res, resolved)
				}
			}
		}
	}

	return res
}

// linkNeighbours populates the neighbours array for each cell
func linkNeighbours(in <-chan []cell) <-chan []cell {
	out := make(chan []cell)
	go func() {
		defer close(out)

		// We buffer an extra row to ensure that if you receive a row, you
		// can fully resolve its neighbours (i.e., if you find a number below
		// the row you're processing then it will definitely be linked).
		var last []cell
		var waiting []cell
		for row := range in {
			for i := range row {
				if last != nil {
					if i > 0 {
						row[i].neighbours[topLeft] = &last[i-1]
						last[i-1].neighbours[bottomRight] = &row[i]

						last[i].neighbours[left] = &last[i-1]
					}

					row[i].neighbours[top] = &last[i]
					last[i].neighbours[bottom] = &row[i]

					if i < len(last)-1 {
						row[i].neighbours[topRight] = &last[i+1]
						last[i+1].neighbours[bottomLeft] = &row[i]

						last[i].neighbours[right] = &last[i+1]
					}
				}
			}

			if waiting != nil {
				out <- waiting
			}

			waiting = last
			last = row
		}

		out <- waiting
		out <- last
	}()
	return out
}

// addBorders adds an empty row of cells above and below the schematic
func addBorders(in <-chan []cell) <-chan []cell {
	out := make(chan []cell)
	go func() {
		defer close(out)
		length := 0
		for row := range in {
			if length == 0 {
				length = len(row)
				out <- make([]cell, length)
			}
			out <- row
		}
		out <- make([]cell, length)
	}()
	return out
}

// parseSchematic takes the raw file and emits rows of individual cells
func parseSchematic(in <-chan byte) <-chan []cell {
	out := make(chan []cell)
	go func() {
		defer close(out)

		var id = 0
		var row []cell
		for b := range in {
			switch {
			case b == '\n' || b == '\r':
				if len(row) > 0 {
					out <- row
					row = make([]cell, 0)
				}
			case b >= '0' && b <= '9':
				row = append(row, cell{id: id, kind: cellKindNumber, value: int(b - '0')})
			case b == '.':
				row = append(row, cell{id: id, kind: cellKindEmpty})
			default:
				row = append(row, cell{id: id, kind: cellKindSymbol, value: int(b)})
			}
			id++
		}
	}()
	return out
}

type cellKind uint8

const (
	cellKindEmpty cellKind = iota
	cellKindNumber
	cellKindSymbol
)

const (
	topLeft     = 0
	left        = 1
	bottomLeft  = 2
	top         = 3
	bottom      = 4
	topRight    = 5
	right       = 6
	bottomRight = 7
)

type cell struct {
	id         int
	neighbours [8]*cell
	kind       cellKind
	value      int
}

func (c *cell) is(kind cellKind) bool {
	return c != nil && c.kind == kind
}

func (c *cell) symbol() byte {
	switch c.kind {
	case cellKindSymbol:
		return byte(c.value)
	case cellKindNumber:
		return byte(c.value + '0')
	default:
		return '.'
	}
}
