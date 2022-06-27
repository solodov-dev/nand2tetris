package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	input := os.Args[1]

	name, dir, isFile := PathInfo(input)

	if !isFile {
		input = dir + "*.vm"
	}

	files, err := filepath.Glob(input)

	if err != nil {
		log.Fatalf("Error reading file")
	}

	writer := NewCodeWriter(dir + name + ".asm")
  writer.WriteInit(isSysFileProvided(files))

	defer writer.Close()

	for _, f := range files {
    // currFileName, _, _ := PathInfo(f)
	  parser := NewParser(f)
		for parser.HasMoreCommands() {
			parser.Advance()
			commandType, err := parser.ParseCommand()

			if err != nil {
				log.Fatalf("Cannot parse command %s on line %d", parser.currentCommand, parser.lineNumber)
			}

			switch commandType {
			case C_PUSH:
				writer.WritePush(parser.arg1, parser.arg2, parser.currentFile)
			case C_ARITHMETIC:
				writer.WriteArithmetic(parser.arg1)
			case C_POP:
				writer.WritePop(parser.arg1, parser.arg2, parser.currentFile)
			case C_LABEL:
				writer.WriteLabel(parser.arg1, parser.currentFunction)
			case C_GOTO:
				writer.WriteGoTo(parser.arg1, parser.currentFunction)
			case C_IF:
				writer.WriteIfGoTo(parser.arg1, parser.currentFunction)
			case C_FUNCTION:
				writer.WriteFunction(parser.arg1, parser.arg2)
			case C_RETURN:
				writer.WriteReturn()
			case C_CALL:
				writer.WriteCall(parser.arg1, parser.arg2)
			}
		}
	}
	writer.WriteEnd()
}

func isSysFileProvided(files []string) bool {
  for _, f := range files {
    if strings.HasSuffix(f, "Sys.vm") {
      return true
    }
  }
  return false
}
