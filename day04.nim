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
      digit: uint8
      ours: bool

    for i in 10..slice.size:
      if line[i] == '|':
        ours = true
      else:
        if line[i] >= '0' and line[i] <= '9':
          digit = 10 * digit + uint8(line[i].ord - '0'.ord)
        elif digit > 0:
          if not ours and digit < 64:
            winningLow.setBit(digit)
          elif not ours:
            winningHigh.setBit(digit - 64)
          elif digit < 64:
            ourLow.setBit(digit)
          else:
            ourHigh.setBit(digit - 64)
            
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
