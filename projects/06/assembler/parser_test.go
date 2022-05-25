package assembler

import (
	"bufio"
	"io"
	"testing"
)

type R struct {
	Data string
	done bool
}

func (r *R) Read(p []byte) (n int, err error) {
	copy(p, []byte(r.Data))
	if r.done {
		return 0, io.EOF
	}
	r.done = true
	return len([]byte(r.Data)), nil
}

func NewR(data string) *R {
	return &R{data, false}
}

func NewTestParser(file string) *Parser {
	reader := NewR(file)
	scanner := bufio.NewScanner(reader)
	return NewParser(scanner)
}

func TestHasMoreCommands(t *testing.T) {
	parser := NewTestParser("LOOP")
	got := parser.HasMoreComands()
	want := true

	if got != want {
		t.Errorf("got %t, wanted %t", got, want)
	}

	got = parser.HasMoreComands()
	want = false

  test(t, got, want)
}

func TestAdvance(t *testing.T) {
	parser := mockTest("(LOOP)\nJGT;LOOP")

	got := parser.currentCommand
	want := "(LOOP)"

  test(t, got, want)

	if parser.HasMoreComands() {
		parser.Advance()
	}

	got = parser.currentCommand
	want = "JGT;LOOP"

  test(t, got, want)
}

func TestParserDest(t *testing.T) {
  parser := mockTest("D=D-A")

  got := parser.Dest()
  want := "D"

  test(t, got, want)
}

func TestParserSymbol(t *testing.T) {
  parser := mockTest("@LOOP")

  got := parser.Symbol()
  want := "LOOP"

  test(t, got, want)
}

func TestParserComp(t *testing.T) {
  parser := mockTest("M=!D")

  got := parser.Comp()
  want := "!D"

  test(t, got, want)
}

func TestParserJump(t *testing.T) {
  parser := mockTest("M=JGT;LOOP")

  got := parser.Jump()
  want := "LOOP"

  test(t, got, want)
}

func mockTest(file string) *Parser {
  parser := NewTestParser(file)

  if parser.HasMoreComands() {
    parser.Advance()
  }

  return parser
}
