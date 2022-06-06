package main

import (
	"bufio"
	"io"
	"testing"
)

type Reader struct {
	Data string
	done bool
}

func (r *Reader) Read(p []byte) (n int, err error) {
	copy(p, []byte(r.Data))
	if r.done {
		return 0, io.EOF
	}
	r.done = true
	return len([]byte(r.Data)), nil
}

func NewReader(data string) *Reader {
	return &Reader{data, false}
}

func NewTestParser(file string) *Parser {
	reader := NewReader(file)
	scanner := bufio.NewScanner(reader)
	return NewParser(scanner)
}

func TestHasMoreCommands(t *testing.T) {
	parser := NewTestParser("LOOP")
	compare(t, parser.HasMoreComands(), true)
	compare(t, parser.HasMoreComands(), false)
}

func TestAdvance(t *testing.T) {
	parser := mockTest("(LOOP)\nJGT;LOOP")
	compare(t, parser.currentCommand, "(LOOP)")

	if parser.HasMoreComands() {
		parser.Advance()
	}

	compare(t, parser.currentCommand, "JGT;LOOP")
}

func TestParserDest(t *testing.T) {
	parser := mockTest("D=D-A")
	compare(t, parser.Dest(), "D")
}

func TestParserSymbol(t *testing.T) {
	parser := mockTest("@LOOP")
	compare(t, parser.Symbol(), "LOOP")
}

func TestParserComp(t *testing.T) {
	parser := mockTest("M=!D")
	compare(t, parser.Comp(), "!D")
}

func TestParserJump(t *testing.T) {
	parser := mockTest("M=JGT;LOOP")
	compare(t, parser.Jump(), "LOOP")
}

func mockTest(file string) *Parser {
	parser := NewTestParser(file)
	if parser.HasMoreComands() {
		parser.Advance()
	}
	return parser
}
