import math, memfiles, sequtils

const identifierSize = uint(sizeof(uint16))

const z: uint16 = 20
const zzz: uint16 = ((z*20+z) shl 5) or z

# Not all letters seem to be used, so we can squish the identifiers a bit to
# save allocating as much memory. Can't get down to 4 bits, though :(
const idMap: seq[uint16] = @[
  0,  # A
  1,  # B
  2,  # C
  3,  # D
  0,  # 
  4,  # F
  5,  # G
  6,  # H
  0,  #
  7,  # J
  8,  # K
  9,  # L
  10, # M
  11, # N
  0,  #
  12, # P
  13, # Q
  14, # R
  15, # S
  16, # T
  0,  # 
  18, # V
  0,  # 
  19, # X
  0,  #
  20, # Z
]

func identifier(s: cstring, offset: int): uint16 {.inline.} =
  # We want the last letter to be in the low 5 bits so we can check
  # it using a mask, but we can compact the other two letters together.
  let
    a = idMap[uint16(s[offset].ord - 'A'.ord)]
    b = idMap[uint16(s[offset+1].ord - 'A'.ord)]
    c = idMap[uint16(s[offset+2].ord - 'A'.ord)]
  result = ((a*20+b) shl 5) or c

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
  let map = alloc(zzz*identifierSize*2)
 
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
  echo countSteps(map, 0, zzz, directions)
  # Part two: run each *A node until we get to a *Z node, then LCM the results
  # This happens to work on all inputs for $reasons, despite not being specified or hinted at...
  echo math.lcm(starts.mapIt(countSteps(map, it, z, directions)))

main()
