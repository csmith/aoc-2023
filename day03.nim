import memfiles

proc digit(ch: char): bool {.inline.} =
  return ch >= '0' and ch <= '9'

proc symbol(ch: char): bool {.inline.} =
  let c = ch.ord
  return (c >= 33 and c <= 45) or c == 47 or c >= 58

proc main(): void =
  let input = memfiles.open("inputs/03.txt")
  
  var lineLength = 0
  var partOne = 0
  var partTwo = 0
  
  var lines = newSeq[cstring](3)

  proc checkSymbol(i: int): bool =
    return (i > 0 and (symbol(lines[0][i-1]) or symbol(lines[1][i-1]) or symbol(lines[2][i-1]))) or 
      (symbol(lines[0][i]) or symbol(lines[1][i]) or symbol(lines[2][i])) or
      (i < lineLength and (symbol(lines[0][i+1]) or symbol(lines[1][i+1]) or symbol(lines[2][i+1])))

  proc readDigit(line: int, i: int): int =
    var o = i
    while o > 0 and digit(lines[line][o-1]):
      o -= 1

    var digit = 0
    while o < lineLength and digit(lines[line][o]):
      digit = digit * 10 + (lines[line][o].ord - '0'.ord)
      o += 1

    return digit

  proc findGearNumbers(i: int): void =
    let topLeft = i > 0 and digit(lines[0][i-1])
    let topMiddle = digit(lines[0][i])
    let topRight = i < lineLength and digit(lines[0][i+1])
    let middleLeft = i > 0 and digit(lines[1][i-1])
    let middleRight = i < lineLength and digit(lines[1][i+1])
    let bottomLeft = i > 0 and digit(lines[2][i-1])
    let bottomMiddle = digit(lines[2][i])
    let bottomRight = i < lineLength and digit(lines[2][i+1])
    
    let count = int(topLeft) + int((not topLeft) and topMiddle) + int((not topMiddle) and topRight) +
       int(middleLeft) + int(middleRight) +
       int(bottomLeft) + int((not bottomLeft) and bottomMiddle) + int((not bottomMiddle) and bottomRight)

    if count != 2:
      return

    var product = 1
    if topLeft:
      product *= readDigit(0, i-1)
    elif topMiddle:
      product *= readDigit(0, i)
    if topRight and not topMiddle:
      product *= readDigit(0, i+1)
    if middleLeft:
      product *= readDigit(1, i-1)
    if middleRight:
      product *= readDigit(1, i+1)
    if bottomLeft:
      product *= readDigit(2, i-1)
    elif bottomMiddle:
      product *= readDigit(2, i)
    if bottomRight and not bottomMiddle:
      product *= readDigit(2, i+1)
    
    partTwo += product

  proc process(): void =
    var digit = 0
    var symbol = false
    for i in 0..lineLength:
      if digit(lines[1][i]):
        digit = digit * 10 + (lines[1][i].ord - '0'.ord)
        symbol = symbol or checkSymbol(i)
      elif digit > 0:
        if symbol:
          partOne += digit
          symbol = false
        digit = 0
      if lines[1][i] == '*':
        findGearNumbers(i)
  
  proc shift(line: cstring): void =
    lines[0] = lines[1]
    lines[1] = lines[2]
    lines[2] = line

  var first = true
  for slice in memSlices(input):
    if first:
      lineLength = slice.size
      shift(cstring(newString(lineLength)))
      first = false
    shift(cast[cstring](slice.data))
    process()
    
  shift(cstring(newString(lineLength)))
  process()

  echo partOne
  echo partTwo

main()
