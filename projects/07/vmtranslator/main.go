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

	parser := NewParser(input)
	writer := NewCodeWriter(output, parser)
	defer writer.Close()

	for parser.HasMoreCommands() {
		parser.Advance()
		commandType, err := parser.ParseCommand()

		if err != nil {
			log.Fatalf("Cannot parse command %s on line %d", parser.currentCommand, parser.lineNumber)
		}

		switch commandType {
		case C_PUSH:
			writer.WritePush()
		case C_ARITHMETIC:
			writer.WriteArithmetic()
		}
	}
}
