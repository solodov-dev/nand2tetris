package assembler

import "testing"

func test[C comparable](t *testing.T, got C, want C) {
  if got != want {
		t.Errorf("got %v, wanted %v", got, want)
  }
}
