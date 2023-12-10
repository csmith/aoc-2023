import memfiles

proc main(): void =
  let input = memfiles.open("inputs/02.txt")
  var partOne = 0
  var partTwo = 0
  
  var digit = 0
  var game = 0
  
  var minRed = 0
  var minBlue = 0
  var minGreen = 0

  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    for i in 0..slice.size:
      let b = line[i]
      if b >= '0' and b <= '9':
        digit = 10 * digit + (b.ord - '0'.ord)
      elif b == ':':
        game = digit
        digit = 0
      elif b == 'r' and digit > 0:
        minRed = max(minRed, digit)
        digit = 0
      elif b == 'g' and digit > 0:
        minGreen = max(minGreen, digit)
        digit = 0
      elif b == 'b' and digit > 0:
        minBlue = max(minBlue, digit)
        digit = 0
    
    if minRed <= 12 and minGreen <= 13 and minBlue <= 14:
      partOne += game
    partTwo += minRed * minGreen * minBlue
    minRed = 0
    minGreen = 0
    minBlue = 0

  echo partOne
  echo partTwo

main()
