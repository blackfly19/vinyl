package main

func encode(originalFile []byte, updatedFile []byte, blockSize int) int {

	prime := 65521
	hashmap := make(map[int]int)

	for i := 0; i < len(originalFile); i = i + blockSize {
		hashmap[adler32(originalFile[i:i+blockSize], blockSize, prime)] = i
	}

	initialHash := adler32(updatedFile[:blockSize], blockSize, prime)

	for i := blockSize; i < len(updatedFile); i++ {
		if hashmap
	}
}

func rollingHash(hash int, data []byte, blockSize int, prime int, pos int) int {

	s1 := (hash - int(data[pos-blockSize]+data[pos])) % prime
	s2 := (hash - blockSize*int(data[pos-blockSize]) + s1) % prime
	return s2<<16 + s1
}

func adler32(data []byte, blockSize int, prime int) int {
	s1 := 1
	s2 := 0

	for i := 0; i < blockSize; i++ {
		s1 = (s1 + int(data[i])) % prime
		s2 = (s2 + s1) % prime
	}
	return s2<<16 + s1
}
