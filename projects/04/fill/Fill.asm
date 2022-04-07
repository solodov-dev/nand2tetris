// This file is part of www.nand2tetris.org
// and the book "The Elements of Computing Systems"
// by Nisan and Schocken, MIT Press.
// File name: projects/04/Fill.asm

// Runs an infinite loop that listens to the keyboard input.
// When a key is pressed (any key), the program blackens the screen,
// i.e. writes "black" in every pixel;
// the screen should remain fully black as long as the key is pressed. 
// When no key is pressed, the program clears the screen, i.e. writes
// "white" in every pixel;
// the screen should remain fully clear as long as no key is pressed.

// Put your code here.

// Watch for keyboard
(LOOP)
  @KBD
  D=M
  @BLACK
  D;JGT
  @WHITE
  0;JMP

// Set color to black
(BLACK)
  @COLOR
  M=-1
  @FILL
  0;JMP

// Set color to white
(WHITE)
  @COLOR
  M=0
  @FILL
  0;JMP

(FILL)
  // Set starting pixel index to the first pixel of screen
  @SCREEN
  D=A
  @PIXEL_IDX
  M=D
  (DRAW)
    // Get color
    @COLOR
    D=M
    // Set RAM[PIXEL_IDX] to color
    @PIXEL_IDX
    A=M
    M=D
    // Increase pixel index
    @PIXEL_IDX
    MD=M+1
    // Last pixel index is KBD memory address - 1
    // Check that we are still in the screen RAM boundaries
    @KBD
    D=D-A
    @LOOP
    D;JEQ
  @DRAW
  0;JMP




