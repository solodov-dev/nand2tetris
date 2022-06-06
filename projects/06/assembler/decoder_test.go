package main

import "testing"

func TestDest(t *testing.T) {
	testCode(t, destDict, Dest)
}

func TestJump(t *testing.T) {
	testCode(t, jumpDict, Jump)
}

func TestComp(t *testing.T) {
	testCode(t, compDict, Comp)
}

func testCode(t *testing.T, dict map[string]string, f func(string) (string, error)) {
	for cmd, want := range dict {
		got, err := f(cmd)

		if got != want || err != nil {
			t.Errorf("got %s, wanted %s", got, want)
		}
	}
}
