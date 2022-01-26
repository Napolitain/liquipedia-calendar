package main

type uidSeed struct {
	data  []byte
	index int
}

func newUidSeed() *uidSeed {
	return &uidSeed{
		[]byte("0000000000"),
		9,
	}
}

func incrementUidSeed(seed *uidSeed) {
	flag := false
	for seed.data[seed.index] == 'z' {
		flag = true
		seed.index--
	}
	if flag {
		for i := seed.index + 1; i < len(seed.data); i++ {
			seed.data[i] = '0'
		}
	}
	if seed.index == -1 {
		return
	}
	seed.data[seed.index]++
	if flag == true {
		seed.index = 9
	}
}
