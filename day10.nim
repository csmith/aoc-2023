import memfiles

const
  empty: uint8 = 0b0000000
  start: uint8 = 0b0000001
  left: uint8 = 0b0000010
  right: uint8 = 0b0000100
  up: uint8 = 0b0001000
  down: uint8 = 0b0010000

func opposite(direction: uint8): uint8 {.inline.} =
  case direction:
    of left:
      right
    of right:
      left
    of up:
      down
    of down:
      up
    else:
      empty

func pickDirection(options: uint8, exclude: uint8): uint8 {.inline.} =
  result = options and not exclude

func canGo(source: uint8, direction: uint8): bool {.inline.} =
  (source and direction) == direction

func move(position: int, direction: uint8, lineLength: int): int {.inline.} =
  result = position
  if direction.canGo(left):
    result -= 1
  if direction.canGo(right):
    result += 1
  if direction.canGo(up):
    result -= lineLength
  if direction.canGo(down):
    result += lineLength

proc main(): void =
  let input = memfiles.open("inputs/10.txt")

  var map = newSeq[uint8](input.size)
  var lineLength = 0
  var startingPos = 0

  var index = 0
  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    if lineLength == 0:
      lineLength = slice.size
    for i in 0..slice.size-1:
      map[index] = (
        case line[i]:
          of 'S':
            start
          of '|':
            up or down
          of '-':
            left or right
          of 'J':
            up or left
          of 'L':
            up or right
          of '7':
            down or left
          of 'F':
            down or right
          else:
            empty
      )
      if line[i] == 'S':
        startingPos = index
      inc index

  # Add the connections in for the start
  if startingPos >= lineLength and map[move(startingPos, up, lineLength)].canGo(down):
    map[startingPos] = map[startingPos] or up
  if startingPos + lineLength <= input.size and map[move(startingPos, down, lineLength)].canGo(up):
    map[startingPos] = map[startingPos] or down
  if (startingPos mod lineLength) < lineLength-1 and map[move(startingPos, left, lineLength)].canGo(right):
    map[startingPos] = map[startingPos] or left
  if (startingPos mod lineLength) > 0 and map[move(startingPos, right, lineLength)].canGo(left):
    map[startingPos] = map[startingPos] or right

  # Pick a direction to start in
  var
    direction = (
      if map[startingPos].canGo(left):
        left
      elif map[startingPos].canGo(right):
        right
      elif map[startingPos].canGo(up):
        up
      else:
        down
    )
    route = newSeq[bool](input.size)
    moves = 0
    pos = startingPos
  
  while true:
    route[pos] = true
    inc moves
    pos = move(pos, direction, lineLength)
    direction = map[pos].pickDirection(direction.opposite)
    if pos == startingPos:
      break

  # Part two: running across each row, we start outside then toggle that every
  # time we cross an upwards connection.
  var holes = 0
  for i in 0..int(input.size/lineLength)-1:
    var inside = false
    for j in 0..lineLength-1:
      if route[i*lineLength+j]:
        if map[i*lineLength+j].canGo(up):
          inside = not inside
      elif inside:
        inc holes
    
  echo int(moves/2)
  echo holes

main()