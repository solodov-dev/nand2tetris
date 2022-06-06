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
  lineNumber uint16
}

type ParserInterface interface {
	Advance()
	CommandType() Command
  HasMoreComands() bool
	Symbol() string
	Dest() string
	Comp() string
	Jump() string
  CurrentLine() string
  CurrentLineNumber() uint16
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

func (p *Parser) HasMoreComands() bool {
	hasCommands := p.scanner.Scan()

	if !hasCommands {
		return false
	}

	if ignoreNextLine(p) {
		return p.HasMoreComands()
	}

  p.lineNumber++

	return true
}

func (p *Parser) Advance() {
  p.currentCommand = strings.TrimSpace(p.scanner.Text())
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
		return ""
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
  if (endIndex < 0) {
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

func ignoreNextLine(p *Parser) bool {
	nextLine := p.scanner.Text()
	if nextLine == "" || strings.HasPrefix(nextLine, "//") {
		return true
	}
	return false
}
