@256
D=A
@SP
M=D

(SimpleFunction.test)
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
@0
D=A
@SP
A=M
M=D
@SP
M=M+1
@LCL 
D=M
@0
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
@LCL 
D=M
@1
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1
@SP
M=M-1
A=M
M=!M
@SP
M=M+1
@ARG 
D=M
@0
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=D+M
@SP
M=M+1
@ARG 
D=M
@1
A=D+A
D=M
@SP
A=M
M=D
@SP
M=M+1
@SP
M=M-1
A=M
D=M
@SP
M=M-1
A=M
M=M-D
@SP
M=M+1
@LCL  // FRAME=LCL  Save LCL in a temp variable
D=M
@R13
M=D
@5  // RET=*(FRAME-5)  Put the return address in a temp var R14
A=D-A
D=M
@R14
M=D
@SP // *ARG=pop()  Reposition the return value
M=M-1
@ARG
AD=M
@R15
M=D
@SP
A=M
D=M
@R15
A=M
M=D
@R2
D=M
@R0
M=D+1
@R13  // THAT=*(FRAME-1)  Restore THAT of the caller
D=M
D=D-1
@R13
M=D
A=D
D=M
@THAT
M=D
@R13  // THIS=*(FRAME-2)  Restore THIS of the caller
D=M
D=D-1
@R13
M=D
A=D
D=M
@THIS
M=D
@R13  // ARG=*(FRAME-3)  Restore ARG of the caller
D=M
D=D-1
@R13
M=D
A=D
D=M
@ARG
M=D
@R13  // LCL=*(FRAME-4)  Restore LCL of the caller
D=M
D=D-1
@R13
M=D
A=D
D=M
@LCL
M=D
@R14
A=M
0;JMP
(END)
@END
0;JMP
