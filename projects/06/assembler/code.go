package assembler

import (
	"strings"
  "errors"
)

func Dest(command string) string {
	binary := [3]string{"0", "0", "0"}

	if strings.Contains(command, "A") {
		binary[0] = "1"
	}

	if strings.Contains(command, "D") {
		binary[1] = "1"
	}

	if strings.Contains(command, "M") {
		binary[2] = "1"
	}

	return strings.Join(binary[:], "")
}

func Comp(command string) (string, error) {
	switch command {
	// a=0
	case "0":
		return "0101010", nil
	case "1":
		return "0111111", nil
	case "-1":
		return "0111010", nil
	case "D":
		return "0001100", nil
	case "A":
		return "0110000", nil
	case "!D":
		return "0001101", nil
	case "!A":
		return "0110001", nil
	case "-D":
		return "0001111", nil
	case "-A":
		return "0110011", nil
	case "D+1":
		return "0011111", nil
	case "A+1":
		return "0110111", nil
	case "D-1":
		return "0001110", nil
	case "A-1":
		return "0110010", nil
	case "D+A":
		return "0000010", nil
	case "D-A":
		return "0010011", nil
	case "A-D":
		return "0000111", nil
	case "D&A":
		return "0000000", nil
	case "D|A":
		return "0010101", nil

	// a=1
	case "M":
		return "1110000", nil
	case "!M":
		return "1110001", nil
	case "-M":
		return "1110011", nil
	case "M+1":
		return "1110111", nil
	case "M-1":
		return "1110010", nil
	case "D+M":
		return "1000010", nil
	case "D-M":
		return "1010011", nil
	case "M-D":
		return "1000111", nil
	case "D&M":
		return "1000000", nil
	case "D|M":
		return "1010101", nil
	}
	return "", errors.New("undefined comp command: " + command)
}

func Jump(command string) (string, error) {
	switch command {
	case "null":
		return "000", nil
	case "JGT":
		return "001", nil
	case "JEQ":
		return "010", nil
	case "JGE":
		return "011", nil
	case "JLT":
		return "100", nil
	case "JNE":
		return "101", nil
	case "JLE":
		return "110", nil
	case "JMP":
		return "111", nil
	}
  return "", errors.New("undefined jump command: " + command)
}
