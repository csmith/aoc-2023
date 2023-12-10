import memfiles, sequtils

func prediction(series: seq[int]): (int, int) {.inline.} =
  var
    lastDigits = 0
    firstDigits = 0
    current = series
    odd = false
  while not all(current, proc (x: int): bool = x == 0):
    firstDigits += current[0] * (if odd: -1 else: 1)
    lastDigits += current[len(current)-1]
    for i in 0..current.len-2:
      current[i] = current[i+1] - current[i]
    current.setLen(current.len-1)
    odd = not odd
  return (lastDigits, firstDigits)
  

proc main(): void =
  let input = memfiles.open("inputs/09.txt")
  
  var
    series = newSeq[int]()
    digit = 0
    negative = false
    partOne = 0
    partTwo = 0

  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    for i in 0..slice.size:
      let b = line[i]
      if b >= '0' and b <= '9':
        digit = digit * 10 + (b.ord - '0'.ord)
      elif b == '-':
        negative = true
      else:
        if negative:
          digit *= -1
        series.add(digit)
        digit = 0
        negative = false
    let (p1, p2) = prediction(series)
    partOne += p1
    partTwo += p2
    series.setLen(0)

  echo partOne
  echo partTwo

main()
