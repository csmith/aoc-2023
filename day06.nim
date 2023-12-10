import math, memfiles

func calcWaysToWin(time: uint, record: uint): uint {.inline.} =
  let target = float64(record)+0.000001
  let sqrt = math.sqrt(float64(time*time) - 4*target)
  let lower = int(math.floor((float64(time) - sqrt) / 2))
  let upper = int(math.ceil((float64(time) + sqrt) / 2 - 1))
  return uint(upper - lower)

proc main(): void =
  let input = memfiles.open("inputs/06.txt")

  var
    digit: uint = 0
    times: seq[uint]
    distances: seq[uint]
    partTwoDistance: uint
    partTwoTime: uint
    first = true

  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    for i in 0..slice.size:
      if line[i] >= '0' and line[i] <= '9':
        let n = uint(line[i].ord - '0'.ord)
        digit = digit * 10 + n
        if first:
          partTwoTime = partTwoTime*10 + n
        else:
          partTwoDistance = partTwoDistance*10 + n
      elif digit > 0:
        if first:
          times.add(digit)
        else:
          distances.add(digit)
        digit = 0
    first = false
  
  var
    partOne: uint = 1
  
  for i in 0..times.len-1:
    partOne *= calcWaysToWin(times[i], distances[i])

  echo partOne
  echo calcWaysToWin(partTwoTime, partTwoDistance)

main()
