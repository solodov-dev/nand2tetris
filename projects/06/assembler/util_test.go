package main

import "testing"

func compare[C comparable](t *testing.T, got C, want C) {
  if got != want {
		t.Errorf("got %v, wanted %v", got, want)
  }
}
