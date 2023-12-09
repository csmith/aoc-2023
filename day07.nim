import algorithm, memfiles

const
  highCard: uint = 0
  onePair: uint = 1
  twoPair: uint = 2
  threeOfAKind: uint = 3
  fullHouse: uint = 4
  fourOfAKind: uint = 5
  fiveOfAKind: uint = 6

  jokerMutationMap: array[7, seq[uint]] = [
    # Each row corresponds to a hand, and each entry in the row corresponds to the number of jokers
    @[highCard, onePair, threeOfAKind, fourOfAKind, fiveOfAKind, fiveOfAKind],
    @[onePair, threeOfAKind, fourOfAKind, fiveOfAKind],
    @[twoPair, fullHouse],
    @[threeOfAKind, fourOfAKind, fiveOfAKind],
    @[fullHouse],
    @[fourOfAKind, fiveOfAKind],
    @[fiveOfAKind]
  ]

func handScore(dupes: array[5, uint]): uint {.inline.} =
  if dupes[4] == 1:
    return fiveOfAKind
  elif dupes[3] == 1:
    return fourOfAKind
  elif dupes[2] == 1 and dupes[1] == 1:
    return fullHouse
  elif dupes[2] == 1:
    return threeOfAKind
  elif dupes[1] == 2:
    return twoPair
  elif dupes[1] == 1:
    return onePair
  else:
    return highCard

proc main(): void =
  let input = memfiles.open("inputs/07.txt")

  var
    partOneCards = newSeqOfCap[(uint, uint)](1000)
    partTwoCards = newSeqOfCap[(uint, uint)](1000)
    partOneCard: uint = 0
    partTwoCard: uint = 0
    bid: uint = 0

    dupes: array[5, uint]
    counts: array[14, uint]

  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    for i in 0..slice.size:
      let chr = line[i]
      if i < 5:
        let value = case chr:
          of '2'..'9': uint(chr.ord - '1'.ord)
          of 'T': 9
          of 'J': 10
          of 'Q': 11
          of 'K': 12
          of 'A': 13
          else: raise newException(Exception, chr & " is not a valid card value")

        if counts[value] > 0:
          dupes[counts[value]-1] -= 1
        counts[value] += 1
        dupes[counts[value]-1] += 1
        
        partOneCard = (partOneCard shl 4) or value
        if value == 10:
          partTwoCard = partTwoCard shl 4
        else:
          partTwoCard = (partTwoCard shl 4) or value
      elif chr >= '0' and chr <= '9':
        bid = bid * 10 + uint(chr.ord - '0'.ord)
   
    partOneCard = (handScore(dupes) shl 20) or partOneCard
    
    # For part two, we work out the hand without jokers and then map the hand
    # to a new one based on the number of jokers
    let jokers = counts[10]
    if jokers > 0:
      dupes[jokers-1] -= 1

    partTwoCard = (jokerMutationMap[handScore(dupes)][jokers] shl 20) or partTwoCard

    partOneCards.add((partOneCard, bid))
    partTwoCards.add((partTwoCard, bid))
    
    # Reset all the random crap for the next line
    for i in 0..13:
      counts[i] = 0
    for i in 0..4:
      dupes[i] = 0

    partOneCard = 0
    partTwoCard = 0
    bid = 0

  partOneCards.sort do (x, y: (uint, uint)) -> int:
    result = cmp(x[0], y[0])

  partTwoCards.sort do (x, y: (uint, uint)) -> int:
    result = cmp(x[0], y[0])

  var 
    partOne: uint = 0
    partTwo: uint = 0

  for i in 0..len(partOneCards)-1:
    partOne += uint(i + 1) * partOneCards[i][1]
    partTwo += uint(i + 1) * partTwoCards[i][1]

  echo partOne
  echo partTwo

main()