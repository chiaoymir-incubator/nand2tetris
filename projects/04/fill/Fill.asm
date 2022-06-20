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

// Pseudo Code
// Goal: Fill the whole screen if any of the keys is pressed (512 * 256, addr = SCREEN + row * 32 + col / 16 )
// Tips: MAX = 256 * 32 = 8192
// 
// LOOP:
//     addr = 0
//     val = 0
//     if KBD != 0 goto BLACK
// 
// BLACK:
//     val = -1
//     goto RENDER
//     
// RENDER:
//     if addr == 8192 goto LOOP
//     SCREEN[addr] = val
//     addr = addr + 1
//     goto RENDER
// 
// END
// ==========

(LOOP)
    @addr
    M=0
    @val
    M=0

    @KBD
    D=M 
    @BLACK
    D;JNE

    @RENDER
    D;JMP

(BLACK)
    @val 
    M=-1

    @RENDER
    0;JMP

(RENDER)
    @temp
    M=0

    // if addr == 8192 goto LOOP
    @8192
    D=A
    @addr
    D=D-M
    @LOOP
    D;JEQ

    // SCREEN[addr] = val
    // RAM[addr + SCREEN] = val
    @addr 
    D=M 
    @SCREEN
    D=D+A
    @temp
    M=D

    @val
    D=M
    @temp
    A=M 
    M=D

    // addr = addr + 1
    @addr 
    M=M+1

    @RENDER
    0;JMP

(END)
    @END 
    0;JMP