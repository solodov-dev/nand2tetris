package main

import (
	"log"
	"strconv"
)

var BaseAddresses = map[string]string{
	"local":    "@LCL",
	"argument": "@ARG",
	"this":     "@THIS",
	"that":     "@THAT",
	"pointer":  "@R3",
	"temp":     "@R5",
}

// PUSH

func Push(segment string, index string) string {
	seg, ok := BaseAddresses[segment]

	if !ok {
		log.Fatalf("Segment %s is not in known segments list", segment)
	}

	return seg + ` 
D=M
@` + index + `
A=D+A
D=M
` + pushDtoStack()
}


// PushConstant pushes value val onto stack
func PushConstant(val string) string {
	return `@` + val + `
D=A
` + pushDtoStack()
}

func PushStatic(index string, filename string) string {
	return `@` + filename + `.` + index + `
D=M
` + pushDtoStack()
}

func PushTempPointer(val string, addr int) string {
	idx, err := strconv.Atoi(val)

	if err != nil {
		log.Fatalf("Value %s is not an index", val)
	}

	addr += idx

	return `@` + strconv.Itoa(addr) + `
  D=M
  ` + pushDtoStack()
}

// Pushes the contents of D register to stack
func pushDtoStack() string {
	return `@SP
A=M
M=D
@SP
M=M+1`
}

// POP

func Pop(segment string, index string) string {
	seg, ok := BaseAddresses[segment]

	if !ok {
		log.Fatalf("Segment %s is not in known segments list", segment)
	}

	return seg + `
D=M
@` + index + `
D=A+D
@R13
M=D
` + popStackToD() + `
@R13
A=M
M=D`
}

func PopTempPointer(index string, addr int) string {
	idx, err := strconv.Atoi(index)

	if err != nil {
		log.Fatalf("Value %s is not an index", index)
	}

	addr += idx

  return popStackToD() + `
@` + strconv.Itoa(addr) + `
M=D`
}

func PopStatic(index string, filename string) string {
	return popStackToD() + `
@` + filename + `.` + index + `
M=D`
}

// Pops stack and stores the value in D
func popStackToD() string {
	return `@SP
M=M-1
A=M
D=M`
}
