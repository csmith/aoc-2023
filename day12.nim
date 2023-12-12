import math, memfiles, sequtils, tables, threadpool

type Spring = enum
  Broken
  Working
  Unknown

# ==================== Parsing ====================

func parseSprings(line: cstring, i: var int): seq[Spring] {.inline.} =
  while true:
    if line[i] == '#':
      result.add(Spring.Broken)
    elif line[i] == '.':
      result.add(Spring.Working)
    elif line[i] == '?':
      result.add(Spring.Unknown)
    else:
      inc i # Consume the space
      return
    inc i
  
func parseReports(line: cstring, i: var int): seq[int] {.inline.} =
  var digit = 0
  while true:
    if line[i] >= '0' and line[i] <= '9':
      digit = digit * 10 + (line[i].ord - '0'.ord)
    else:
      result.add(digit)
      digit = 0
      if line[i] != ',':
        return
    inc i

# ==================== Memoisation ====================

proc combinations(cache: TableRef[(seq[Spring], seq[int]), int], springs: seq[Spring], reports: seq[int]): int

proc memoisedCombinations(cache: TableRef[(seq[Spring], seq[int]), int], springs: seq[Spring], reports: seq[int]): int {.inline.} =
  result = cache.getOrDefault((springs, reports), -1)
  if result == -1:
    result = combinations(cache, springs, reports)
    cache[(springs, reports)] = result

proc memoisedCombinations(big: bool, springs: seq[Spring], reports: seq[int]): int {.inline.} =
  {.gcsafe.}: # Where we're going we don't need no GC!
    # Create a new cache, as we don't get much benefit from a global one, and it makes threading hard/crashy.
    # For part one problems it can be a smaller cache, for part two pre-allocate it much bigger.
    let cache = newTable[(seq[Spring], seq[int]), int](if big: 1500 else: 100)
    result = memoisedCombinations(cache, springs, reports)

# ==================== Actual logic ===================

proc combinations(cache: TableRef[(seq[Spring], seq[int]), int], springs: seq[Spring], reports: seq[int]): int =
  if unlikely(len(springs) == 0 and len(reports) == 0):
    # Nothing left, all reports have been satisfied, so we're done.
    return 1

  if unlikely(len(reports) == 0):
    # We're not expecting any more broken springs.
    for s in springs:
      if s == Broken:
        return 0
    return 1

  if unlikely(len(springs) < sum(reports) + len(reports) - 1):
    # Not enough springs to satisfy all reports.
    return 0

  if springs[0] == Broken:
    # If the first spring is broken, then we have to satisfy the first report.
    if likely(reports[0] > 1):
      for i in 1..reports[0]-1:
        if springs[i] == Working:
          return 0
    
    if len(springs) > reports[0]:
      # If there are further springs, the next one can't be broken to leave a gap.
      if springs[reports[0]] == Broken:
        return 0
      return memoisedCombinations(cache, springs[(reports[0] + 1)..len(springs)-1], reports[1..len(reports)-1])
    else:
      return 1
  elif springs[0] == Working:
    # Just ignore all the working springs.
    var i = 1
    while i < len(springs) and springs[i] == Working:
      inc i
    return memoisedCombinations(cache, springs[i..len(springs)-1], reports)
  else:
    # First spring is unknown, try both ways.
    return memoisedCombinations(cache, concat(@[Broken], springs[1..len(springs)-1]), reports) +
       memoisedCombinations(cache, springs[1..len(springs)-1], reports)

# ==================== Plumbing ====================

proc main(): void =
  let input = memfiles.open("inputs/12.txt")
  var partOneFlows = newSeqOfCap[FlowVar[int]](1000)
  var partTwoFlows = newSeqOfCap[FlowVar[int]](1000)

  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    var i = 0
    let springs = parseSprings(line, i)
    let reports = parseReports(line, i)
    
    partOneFlows.add(spawn memoisedCombinations(false, springs, reports))
    partTwoFlows.add(spawn memoisedCombinations(true, concat(springs, @[Unknown], springs, @[Unknown], springs, @[Unknown], springs, @[Unknown], springs), cycle(reports, 5)))

  var
    partOne = 0
    partTwo = 0
  for f in partOneFlows:
    partOne += ^f
  for f in partTwoFlows:
    partTwo += ^f
  
  echo partOne
  echo partTwo

main()
