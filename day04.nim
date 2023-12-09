import bitops, memfiles

proc main(): void =
  let input = memfiles.open("inputs/04.txt")

  var
    partOne, partTwo: int
    copies: array[11, int] # We can win 10 times on one card, plus one for the current card
    copiesIdx: int
  
  for i in 0..10:
    copies[i] = 1

  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    
    var
      winningLow, winningHigh, ourLow, ourHigh: uint64  # WTB: uint128
      phase, digit: uint8

    for i in 0..slice.size:
      if line[i] == ':' or line[i] == '|':
        phase += 1
        continue
      elif phase > 0:
        if line[i] >= '0' and line[i] <= '9':
          digit = 10 * digit + uint8(line[i].ord - '0'.ord)
        elif digit > 0:
          if phase == 1 and digit < 64:
            winningLow = winningLow or uint64(1 shl digit)
          elif phase == 1:
            winningHigh = winningHigh or uint64(1 shl (digit - 64))
          elif phase == 2 and digit < 64:
            ourLow = ourLow or uint64(1 shl digit)
          elif phase == 2:
            ourHigh = ourHigh or uint64(1 shl (digit - 64))
            
          digit = 0
    
    partTwo += int(copies[copiesIdx])

    let wins = countSetBits(winningLow and ourLow) + countSetBits(winningHigh and ourHigh)
    if wins > 0:    
      partOne += 1 shl (wins-1)

      for i in 0..(wins-1):
        copies[(copiesIdx + 1 + i) mod 11] += copies[copiesIdx]
    
    copies[copiesIdx] = 1
    copiesIdx = (copiesIdx + 1) mod 11
    
  echo partOne
  echo partTwo

main()