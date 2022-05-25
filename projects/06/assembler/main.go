package main

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	input, output := readArgs()

	readFile := openFile(input)
  defer readFile.Close();

	newFile := createFile(output)

	scanner := bufio.NewScanner(readFile)
	parser := NewParser(scanner)

	for parser.HasMoreComands() {
		parser.Advance()
		switch parser.CommandType() {
		case A_COMMAND:
			symbol := parser.Symbol()
			number, err := strconv.Atoi(symbol)

			if err != nil {
				log.Fatal("Symbol is not a number. Aborting")
			}

			binary := strconv.FormatInt(int64(number), 2)
      padding := strings.Repeat("0", 16 - len(binary))
			_, err = newFile.WriteString(padding + binary + "\n")

			if err != nil {
				log.Fatalf("Could not write line %s to a new file %s", binary, newFile.Name())
			}
		case C_COMMAND:
			dest := Dest(parser.Dest())
			comp, err := Comp(parser.Comp())

      if err != nil {
        log.Fatalf("Cannot read COMP part of command: %s", err)
      }

			jump, err := Jump(parser.Jump())

      if err != nil {
        log.Fatalf("Cannot read JUMP part of command: %s", err)
      }

			binary := "111" + comp + dest + jump

      _, err = newFile.WriteString(binary + "\n")
			if err != nil {
				log.Fatalf("Could not write line %s to a new file %s", binary, newFile.Name())
			}
		case L_COMMAND:
			log.Fatal("Unknown L command")
		}
	}
}

func readArgs() (string, string) {
	args := os.Args[1:]
	if len(args) < 2 {
		log.Fatal("Not enough arguments. Usage: assembler inputfile outputfile")
	}

	input := os.Args[1]
	if input == "" {
		log.Fatal("No input file specified")
	}

	output := os.Args[2]
	if output == "" {
		log.Fatal("No output file specified")
	}

	return input, output
}

func openFile(input string) *os.File {
	file, err := os.Open(input)

	if err != nil {
		log.Fatalf("Cannot open file %s", input)
	}

	return file
}

func createFile(output string) *os.File {
	newFile, err := os.Create(output)

	if err != nil {
		log.Fatalf("Cannot create new file %s", output)
	}

	return newFile
}
