package main

import "testing"

func TestDummy(t *testing.T) {
	if 2+2 != 4 {
		t.Fatal("Math is broken")
	}
}
