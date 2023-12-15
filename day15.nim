import memfiles, tables

proc main(): void =
  let input = memfiles.open("inputs/15.txt")

  var
    partOne = 0
    partTwo = 0
    hash = 0
    label = 0
    labelHash = 0
    value = 0
    lenses = newTable[int, int](600)
    hashmap = newSeq[seq[int]](256)

  proc process() =
    partOne = partOne + hash
    if value > 0:
      lenses[label] = value
      let pos = hashmap[labelHash].find(label)
      if pos == -1:
        hashmap[labelHash].add(label)
    label = 0
    hash = 0
    value = 0

  for slice in memSlices(input):
    let line = cast[cstring](slice.data)
    for i in 0..slice.size-1:
      if line[i] == ',':
        process()
      else:
        if line[i] >= 'a' and line[i] <= 'z':
          label = label * 26 + (line[i].ord - 'a'.ord)
        elif line[i] == '-':
          let pos = hashmap[hash].find(label)
          if pos > -1:
            hashmap[hash].delete(pos)
        elif line[i] == '=':
          labelHash = hash
        elif line[i] >= '0' and line[i] <= '9':
          value = 10 * value + (line[i].ord - '0'.ord)
        hash = (17 * (hash + line[i].ord)) mod 256

    process()
  
  echo partOne

  for i in 0..255:
    for j in 0..(hashmap[i].len-1):
      partTwo += (i+1) * (j+1) * lenses[hashmap[i][j]]

  echo partTwo


main()