package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("Not enough arguments. Usage: assembler inputfile outputfile")
	}

	input := os.Args[1]
	output := os.Args[2]

	symbolTable := NewSymbolTable()

	io := NewIO(input, output)
	defer io.Close()

	parser := NewParser(io)
	decoder := NewDecoder(parser, symbolTable)

	// First pass
	ROMLineNumber := 0
	for parser.HasMoreComands() {
		parser.Advance()
		if parser.CommandType() == L_COMMAND {
			symbol := parser.Symbol()
			if symbolTable.Contains(symbol) {
				log.Fatalf("Double Label Command %s encountered on line %d", symbol, ROMLineNumber)
			}
			symbolTable.AddEntry(symbol, ROMLineNumber)
		} else {
			ROMLineNumber++
		}
	}

	parser.Reset(io)

	// Second pass
	for parser.HasMoreComands() {
		parser.Advance()
		switch parser.CommandType() {
		case A_COMMAND:
			aCommand := decoder.DecodeACommand()
			io.Write(aCommand)
		case C_COMMAND:
			cCommand := decoder.DecodeCCommand()
			io.Write(cCommand)
		case L_COMMAND:
			continue
		}
	}
}
