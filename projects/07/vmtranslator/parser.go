package main

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Command int

const (
	C_ARITHMETIC Command = iota
	C_PUSH
	C_POP
	C_LABEL
	C_GOTO
	C_IF
	C_FUNCTION
	C_RETURN
	C_CALL
	C_UNKNOWN
)

var patterns = map[Command]*regexp.Regexp{
	C_ARITHMETIC: regexp.MustCompile(`^(add|sub|neg|eq|gt|lt|and|or|not)$`),
	C_POP:        regexp.MustCompile(`^pop (\w+) (\w+)$`),
	C_PUSH:       regexp.MustCompile(`^push (\w+) (\w+)$`),
	C_GOTO:       regexp.MustCompile(`^goto (\w+)$`),
	C_FUNCTION:   regexp.MustCompile(`^function ([\w.]+) (\w+)$`),
	C_IF:         regexp.MustCompile(`^if-goto (\w+)$`),
	C_LABEL:      regexp.MustCompile(`^label (\w+)$`),
	C_RETURN:     regexp.MustCompile(`^return$`),
	C_CALL:       regexp.MustCompile(`^call ([\w.]+) (\w+)$`),
}

type Parser struct {
	scanner         *bufio.Scanner
	currentCommand  string
	currentFile     string
	currentFunction string
	lineNumber      uint16
	arg1            string
	arg2            string
}

func NewParser(input string) *Parser {
	readFile, err := os.Open(input)

	if err != nil {
		log.Fatalf("Cannot open file %s", input)
	}

	scanner := bufio.NewScanner(readFile)
	basename := filepath.Base(input)
	nameWithoutSuffix := strings.TrimSuffix(basename, filepath.Ext(basename))

	return &Parser{scanner, "", nameWithoutSuffix, "", 0, "", ""}
}

func (p *Parser) HasMoreCommands() bool {
	hasCommands := p.scanner.Scan()

	if !hasCommands {
		return false
	}

	nextLine := p.scanner.Text()
	p.lineNumber++
	if nextLine == "" || strings.HasPrefix(nextLine, "//") {
		return p.HasMoreCommands()
	}

	return true
}

func (p *Parser) Advance() {
	text := p.scanner.Text()
	index := strings.Index(text, "//")
	if index > -1 {
		text = text[:index]
	}
	text = strings.TrimSpace(text)
	p.currentCommand = text
}

func (p *Parser) ParseCommand() (Command, error) {
	for commandType, regExp := range patterns {
		if match := regExp.FindStringSubmatch(p.currentCommand); len(match) > 0 {
			if len(match) > 1 {
				p.arg1 = match[1]

				if len(match) > 2 {
					p.arg2 = match[2]
				} else {
					p.arg2 = ""
				}
			}

			// Save current function name
			if commandType == C_FUNCTION {
				p.currentFunction = p.arg1
			}

			return commandType, nil
		}
	}

	return C_UNKNOWN, errors.New("Unknown command type")
}

func (p *Parser) Arg1() string {
	return p.arg1
}

func (p *Parser) Arg2() string {
	return p.arg2
}
