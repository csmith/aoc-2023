import memfiles

func matches(line: cstring, offset: int, substr: cstring): bool {.inline.} =
  var i = 0
  while true:
    if substr[i] == '\0': return true
    if substr[i] != line[offset+i]: return false
    inc(i)

proc main(): void =
  let input = memfiles.open("inputs/01.txt")

  const words: seq[cstring] = @["one", "two", "three", "four", "five", "six", "seven", "eight", "nine"]

  var partOne = 0
  var partTwo = 0

  for slice in memSlices(input):
    let line = cast[cstring](slice.data)

    block forward:
      var foundPartTwo = false
      for i in 0..slice.size:
        if line[i] >= '0' and line[i] <= '9':
          partOne += 10 * (line[i].ord - '0'.ord)
          if not foundPartTwo:
            partTwo += 10 * (line[i].ord - '0'.ord)
          break forward
        elif not foundPartTwo:
          for v, w in words:
            if matches(line, i, w):
              foundPartTwo = true
              partTwo += 10 * (v + 1)

    block backward:
      var foundPartTwo = false
      for i in countdown(slice.size-1, 0):
        if line[i] >= '0' and line[i] <= '9':
          partOne += line[i].ord - '0'.ord
          if not foundPartTwo:
            partTwo += line[i].ord - '0'.ord
          break backward
        elif not foundPartTwo:
          for v, w in words:
            if matches(line, i, w):
              foundPartTwo = true
              partTwo += v + 1

  echo partOne
  echo partTwo

main()