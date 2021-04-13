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

(BEGIN)

(LOOP)

@KBD
D=M
@IFUNFILL
D;JEQ
@IFFILL
0;JMP

(IFFILL)
@R0
M=-1
@ENDIF
0;JMP

(IFUNFILL)
@R0
M=0
@ENDIF
0;JMP

(ENDIF)

@R1
M=0
(FILLLOOP)
@R1
D=M
@SCREEN
D=A+D
@R2
M=D
@R0
D=M
@R2
A=M
M=D

@R1
MD=M+1
@8192
D=D-A
@FILLLOOP
D;JNE

@LOOP
0;JMP

(END)