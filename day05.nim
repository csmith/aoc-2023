import memfiles

proc main(): void =
  let input = memfiles.open("inputs/05.txt")
  
  var
    seeds: seq[uint]
    ranges: seq[array[2, uint]]
    mappings: seq[seq[array[3, uint]]]
    first = true

  # Parsing:
  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    
    if first:
      for i in 6..slice.size:  
        if line[i] == ' ':
          seeds.add(0)
        if line[i] >= '0' and line[i] <= '9':
          seeds[len(seeds)-1] = seeds[len(seeds)-1] * 10 + uint(line[i].ord - '0'.ord)  
      first = false
      for i in 0..int((len(seeds)-1)/2):
        ranges.add([seeds[i*2], seeds[i*2+1]])
    elif line[0] >= '0' and line[0] <= '9':
      var
        parts: array[3, uint] # destination, source, length
        part = 0
      for i in 0..slice.size:
        if line[i] >= '0' and line[i] <= '9':
          parts[part] = parts[part] * 10 + uint(line[i].ord - '0'.ord)
        elif line[i] == ' ':
          part += 1
      mappings[len(mappings)-1].add(parts)
    elif len(mappings) == 0 or len(mappings[len(mappings)-1]) > 0:
      mappings.add(@[])

  # Part one:
  for i in 0..len(seeds)-1:
    for m in mappings:
      for row in m:
        if seeds[i] >= row[1] and seeds[i] < row[1] + row[2]:
          seeds[i] += row[0] - row[1]
          break
  let partOne = min(seeds)

  # Part two:
  var newRanges: seq[array[2, uint]]
  var rangeIdx = 0
  
  for m in mappings:   
    while rangeIdx < len(ranges):
      let r = ranges[rangeIdx]
      rangeIdx += 1

      var processed = false
    
      for map in m:
        if map[1] <= r[0] and map[1] + map[2] >= r[0] + r[1]:
          # Seeds are fully contained within map. Simples.
          newRanges.add([r[0] + map[0] - map[1], r[1]])
          processed = true
          break
        elif map[1] > r[0] and map[1] < r[0] + r[1]:
          # The mapping starts in the middle of our seed range, cut it into two
          ranges.add([r[0], map[1] - r[0] - 1])
          ranges.add([map[1], r[1] - (map[1] - r[0])])
          processed = true
          break
        elif map[1] + map[2] > r[0] and map[1] + map[2] < r[0] + r[1]:
          # The mapping ends in the middle of our seed range, cut it into two
          ranges.add([r[0], map[1] + map[2] - r[0] - 1])
          ranges.add([map[1] + map[2], r[1] - (map[1] + map[2] - r[0])])
          processed = true
          break
      
      # Didn't match any of the mappings, keep it as-is
      if not processed:
        newRanges.add(r)

    ranges = newRanges
    rangeIdx = 0
    newRanges.setLen(0)

  var partTwo = high(uint)
  for r in ranges:
    if r[0] < partTwo:
      partTwo = r[0]

  echo partOne
  echo partTwo

main()  