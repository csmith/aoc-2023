import memfiles

func distanceBetween(a: (int, int), b: (int, int), horizontalExpansion: seq[int], verticalExpansion: seq[int], expansionFactor: int): int {.inline.} =
  let x1 = a[0] + (expansionFactor - 1) * horizontalExpansion[a[0]]
  let x2 = b[0] + (expansionFactor - 1) * horizontalExpansion[b[0]]
  let y1 = a[1] + (expansionFactor - 1) * verticalExpansion[a[1]]
  let y2 = b[1] + (expansionFactor - 1) * verticalExpansion[b[1]]
  return abs(x1 - x2) + abs(y1 - y2)

proc main(): void =
  let input = memfiles.open("inputs/11.txt")

  var galaxies = newSeqOfCap[(int, int)](500)
  var columns = newSeq[bool](140)
  var verticalExpansion = newSeq[int](140)
  var y = 0
  var expansionCounter = 0

  # Parse and keep a tally of vertical expansion
  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    var found = false
    for i in 0..slice.size-1:
      if line[i] == '#':
        columns[i] = true
        galaxies.add((i, y))
        found = true
      discard
    if not found:
      inc expansionCounter
    verticalExpansion[y] = expansionCounter
    inc y

  # Work out horizontal expansion
  var horizontalExpansion = newSeq[int](len(columns))
  expansionCounter = 0
  for i in 0..columns.len-1:
    if not columns[i]:
      inc expansionCounter
    horizontalExpansion[i] = expansionCounter

  # Calculate the manhatten distances between everything
  var partOne = 0
  var partTwo = 0
  for i in 0..galaxies.len-1:
    for j in (i+1)..galaxies.len-1:
      partOne += distanceBetween(galaxies[i], galaxies[j], horizontalExpansion, verticalExpansion, 2)
      partTwo += distanceBetween(galaxies[i], galaxies[j], horizontalExpansion, verticalExpansion, 1000000)

  echo partOne
  echo partTwo

main()