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
//     if KBD != 0 goto BLACK_LOOP
// 
// WHITE_LOOP:
//     if addr == 8192 goto LOOP
//     SCREEN[addr] = 0
//     addr = addr + 1
//     goto WHITE_LOOP
// 
// BLACK_LOOP:
//     if addr == 8192 goto LOOP
//     SCREEN[addr] = -1
//     addr = addr + 1
//     goto BLACK_LOOP
// 
// END
// ==========

(LOOP)
    @addr
    M=0

    @KBD
    D=M 
    @BLACK_LOOP
    D;JNE

(WHITE_LOOP)
    // if addr == 8192 goto LOOP
    @8192
    D=A
    @addr
    D=D-M
    @LOOP
    D;JEQ

    // SCREEN[addr] = 0
    @addr 
    D=M
    @SCREEN
    A=A+D
    M=0

    // addr = addr + 1
    @addr 
    M=M+1

    @WHITE_LOOP
    0;JMP

(BLACK_LOOP)
    // if addr == 8192 goto LOOP
    @8192
    D=A
    @addr
    D=D-M
    @LOOP
    D;JEQ

    // SCREEN[addr] = -1
    @addr 
    D=M
    @SCREEN
    A=A+D
    M=-1

    // addr = addr + 1
    @addr 
    M=M+1

    @BLACK_LOOP
    0;JMP

(END)
    @END 
    0;JMP