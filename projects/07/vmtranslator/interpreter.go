package main

// x + y
func add() string {
	return twoArgCommand("M=D+M")
}

// x - y
func sub() string {
	return twoArgCommand("M=M-D")
}

// -y
func neg() string {
	return oneArgCommand("M=-M")
}

// x = y
// true -1
// false 0
func eq() string {
  return compareCommand("JEQ")
}

// if x > y
// true -1
// false 0
func gt() string {
  return compareCommand("JGT")
}

// if x < y
// true -1
// false 0
func lt() string {
  return compareCommand("JLT")
}

// x and y
func and() string {
 return twoArgCommand("M=M&D")
}

// x or y
func or() string {
 return twoArgCommand("M=M|D")
}

// not y
func not() string {
 return oneArgCommand("M=!M")
}

// Create a two arg command on the stack
// Puts X in M register and Y in D register
// Stack pointer points at X
//
//         -----------
//  SP ->  |    x    |    in M
//         -----------
//         |    y    |    in D
//         -----------
// After command is added stack pointer will point at X + 1
func twoArgCommand(cmd string) string {
	return `@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
` + cmd + `
@SP
M=M+1`
}

// Create a one arg command on the stack
// Puts Y in M register
// Stack pointer points at Y
//
//         -----------
//         |    x    |
//         -----------
//  SP ->  |    y    |    in M
//         -----------
// After command is added stack pointer will point at X + 1
func oneArgCommand(cmd string) string {
	return `@SP
M=M-1
A=M
` + cmd + `
@SP
M=M+1`
}

// Compare the last stack value with the compare command
// If false write 0 to M
// If true write -1 to M
func compareCommand(compare string) string {
	return `@SP
M=M-1
A=M
D=M
@TRUE
D;` + compare + `
(FALSE)
D=0
@WRITE
0;JMP
(TRUE)
D=-1
(WRITE)
@SP
A=M
M=D
@SP
M=M+1`
}

var Arithmetic = map[string]func() string{
	"add": add,
	"sub": sub,
	"neg": neg,
  "eq": eq,
  "lt": lt,
  "gt": gt,
  "and": and,
  "or": or,
  "not": not,
}

// Push pushes value val onto stack
func Push(val string) string {
 return `@` + val + `
D=A
@SP
A=M
M=D
@SP
M=M+1`
}
