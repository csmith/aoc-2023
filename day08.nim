import math, memfiles, sequtils

const identifierSize = uint(sizeof(uint16))

func identifier(s: cstring, offset: int): uint16 {.inline.} =
  for i in 0..2:
    result = (result shl 5) or uint16(s[offset+i].ord - 'A'.ord)

func writeDirections(map: pointer, source: uint16, left: uint16, right: uint16): void {.inline.} =
  let a = cast[uint](map) + source*2*identifierSize
  cast[ptr uint16](a)[] = left
  cast[ptr uint16](a+identifierSize)[] = right

func getDirection(map: pointer, source: uint16, right: bool): uint16 {.inline.} =
  let a = cast[uint](map) + source*2*identifierSize + (if right: identifierSize else: 0)
  result = cast[ptr uint16](a)[]

func countSteps(map: pointer, source: uint16, targetMask: uint16, directions: seq[bool]): int {.inline.} =
  var
    directionIdx = 0
    current = source
  while (current and targetMask) != targetMask:
    current = getDirection(map, current, directions[directionIdx])
    directionIdx = (directionIdx + 1) mod directions.len
    inc result

proc main(): void =
  let input = memfiles.open("inputs/08.txt")
  let map = alloc(identifierSize*26425*2)
 
  var first = true
  var directions: seq[bool]
  var starts: seq[uint16]

  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    if first:
      first = false
      directions = newSeq[bool](slice.size)
      for i in 0..slice.size-1:
        directions[i] = line[i] == 'R'
    elif slice.size != 0:
      let source = identifier(line, 0)
      writeDirections(map, source, identifier(line, 7), identifier(line, 12))
      if (source and 0b11111) == 0:
        starts.add(source)

  # Part one: run until we get to a ZZZ node
  echo countSteps(map, 0, 0b110011100111001, directions)
  # Part two: run each *A node until we get to a *Z node, then LCM the results
  # This happens to work on all inputs for $reasons, despite not being specified or hinted at...
  echo math.lcm(starts.mapIt(countSteps(map, it, 0b11001, directions)))

main()
