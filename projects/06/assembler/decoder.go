package main

import (
	"errors"
	"log"
	"strconv"
	"strings"
)

var destDict = map[string]string{
	"null": "000",
	"M":    "001",
	"D":    "010",
	"MD":   "011",
	"A":    "100",
	"AM":   "101",
	"AD":   "110",
	"AMD":  "111",
}

var compDict = map[string]string{
	"0":   "0101010",
	"1":   "0111111",
	"-1":  "0111010",
	"D":   "0001100",
	"A":   "0110000",
	"!D":  "0001101",
	"!A":  "0110001",
	"-D":  "0001111",
	"-A":  "0110011",
	"D+1": "0011111",
	"A+1": "0110111",
	"D-1": "0001110",
	"A-1": "0110010",
	"D+A": "0000010",
	"D-A": "0010011",
	"A-D": "0000111",
	"D&A": "0000000",
	"D|A": "0010101",
	"M":   "1110000",
	"!M":  "1110001",
	"-M":  "1110011",
	"M+1": "1110111",
	"M-1": "1110010",
	"D+M": "1000010",
	"D-M": "1010011",
	"M-D": "1000111",
	"D&M": "1000000",
	"D|M": "1010101",
}

var jumpDict = map[string]string{
	"null": "000",
	"JGT":  "001",
	"JEQ":  "010",
	"JGE":  "011",
	"JLT":  "100",
	"JNE":  "101",
	"JLE":  "110",
	"JMP":  "111",
}

type Decoder struct {
	parser *Parser
}

func NewDecoder(p *Parser) *Decoder {
	return &Decoder{p}
}

func (d *Decoder) DecodeCCommand() string {
	comp, cerr := toBinary(d.parser.Comp(), compDict)
	dest, derr := toBinary(d.parser.Dest(), destDict)
	jump, jerr := toBinary(d.parser.Jump(), jumpDict)

	if cerr != nil || derr != nil || jerr != nil {
		log.Fatalf("Cannot read C command %s on line %d", d.parser.CurrentLine(), d.parser.CurrentLineNumber())
	}

	return "111" + comp + dest + jump
}

func (d *Decoder) DecodeACommand() string {
	symbol := d.parser.Symbol()
	number, err := strconv.Atoi(symbol)

	if err != nil {
		log.Fatal("Symbol is not a number. Aborting")
	}

	binary := strconv.FormatInt(int64(number), 2)
	padding := strings.Repeat("0", 16-len(binary))

	return padding + binary
}

func toBinary(command string, dict map[string]string) (string, error) {
	if val, ok := dict[command]; ok {
		return val, nil
	}

	return "", errors.New("Undefined code command: " + command)
}
