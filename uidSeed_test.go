package main

import (
	"testing"
)

func Test_newUidSeed(t *testing.T) {
	uidSeed := newUidSeed()
	if string(uidSeed.data[:]) != "0000000000" {
		t.Fatal("Value not correctly initialized.")
	}
}

func Test_incrementUidSeed(t *testing.T) {
	uidSeed := newUidSeed()
	for i := 0; i < 5625; i++ {
		incrementUidSeed(uidSeed)
	}
	if string(uidSeed.data[:]) != "0000000100" {
		t.Fatal("Value is incorrectly incremented: " + string(uidSeed.data[:]))
	}
}
