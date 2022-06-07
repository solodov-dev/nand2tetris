package main

import "errors"

type SymbolTable struct {
	table map[string]int
	addr  int
}

var table = map[string]int{
	"SP":     0,
	"LCL":    1,
	"ARG":    2,
	"THIS":   3,
	"THAT":   4,
	"R0":     0,
	"R1":     1,
	"R2":     2,
	"R3":     3,
	"R4":     4,
	"R5":     5,
	"R6":     6,
	"R7":     7,
	"R8":     8,
	"R9":     9,
	"R10":    10,
	"R11":    11,
	"R12":    12,
	"R13":    13,
	"R14":    14,
	"R15":    15,
	"SCREEN": 16384,
	"KBD":    24576,
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{table, 15}
}

func (t *SymbolTable) AddEntry(symbol string, entry int) (int) {
	t.table[symbol] = entry
  return entry
}

func (t *SymbolTable) NextAddress() (int) {
  t.addr++
  return t.addr
}

func (t *SymbolTable) Contains(symbol string) bool {
	_, ok := t.table[symbol]
	return ok
}

func (t *SymbolTable) GetAddress(symbol string) (int, error) {
	if t.Contains(symbol) {
		return t.table[symbol], nil
	}

	return 0, errors.New("No address associated with the symbol " + symbol)
}
