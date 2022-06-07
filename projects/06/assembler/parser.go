package main

import (
	"bufio"
	"strings"
)

type Command int

const (
	A_COMMAND Command = 0
	C_COMMAND Command = 1
	L_COMMAND Command = 2
)

type Parser struct {
	scanner        *bufio.Scanner
	currentCommand string
	lineNumber     uint16
}

func NewParser(io *IO) *Parser {
	scanner := bufio.NewScanner(io.readFile)
	return &Parser{scanner, "", 0}
}

func (p *Parser) CurrentLineNumber() uint16 {
	return p.lineNumber
}

func (p *Parser) CurrentLine() string {
	return p.currentCommand
}

func (p *Parser) Reset(io *IO) {
	io.readFile.Seek(0, 0)
  p.lineNumber = 0
  p.currentCommand = ""
	p.scanner = bufio.NewScanner(io.readFile)
}

func (p *Parser) HasMoreComands() bool {
	hasCommands := p.scanner.Scan()

	if !hasCommands {
		return false
	}

	nextLine := p.scanner.Text()
	p.lineNumber++
	if nextLine == "" || strings.HasPrefix(nextLine, "//") {
		return p.HasMoreComands()
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

func (p *Parser) CommandType() Command {
	if strings.HasPrefix(p.currentCommand, "@") {
		return A_COMMAND
	}

	if strings.HasPrefix(p.currentCommand, "(") {
		return L_COMMAND
	}

	return C_COMMAND
}

func (p *Parser) Symbol() string {
	if p.CommandType() == A_COMMAND {
		return strings.TrimPrefix(p.currentCommand, "@")
	} else {
		return strings.Trim(p.currentCommand, "()")
	}
}

func (p *Parser) Dest() string {
	index := strings.IndexByte(p.currentCommand, '=')
	if index < 0 {
		return "null"
	}
	return p.currentCommand[:index]
}

func (p *Parser) Comp() string {
	startIndex := strings.IndexByte(p.currentCommand, '=')
	if startIndex < 0 {
		startIndex = 0
	} else {
		startIndex++
	}

	endIndex := strings.IndexByte(p.currentCommand, ';')
	if endIndex < 0 {
		return p.currentCommand[startIndex:]
	}
	return p.currentCommand[startIndex:endIndex]
}

func (p *Parser) Jump() string {
	index := strings.IndexByte(p.currentCommand, ';')

	if index < 0 {
		return "null"
	} else {
		index++
	}

	return p.currentCommand[index:]
}
