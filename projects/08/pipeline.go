package main

import (
	"fmt"
	"strconv"
)

type Pipeline struct {
	codes []string
}

func NewPipeline() Pipeline {
	return Pipeline{
		codes: []string{},
	}
}

func (p *Pipeline) pop1() *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%s", SymbolSP),
		"A=M-1",
	)
	return p
}

func (p *Pipeline) pop2() *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%s", SymbolSP),
		"AM=M-1",
		"D=M",
		"A=A-1",
	)
	return p
}

func (p *Pipeline) arithBin(op string) *Pipeline {
	var code string
	switch op {
	case "add":
		code = "M=M+D"
	case "sub":
		code = "M=M-D"
	case "and":
		code = "M=D&M"
	case "or":
		code = "M=D|M"

	}
	p.codes = append(p.codes, code)
	return p
}

func (p *Pipeline) arithUni(op string) *Pipeline {
	var code string
	switch op {
	case "neg":
		code = "M=-M"
	case "not":
		code = "M=!M"
	}
	p.codes = append(p.codes, code)
	return p
}

func (p *Pipeline) arithComp(op string, i int) *Pipeline {
	symbol1 := fmt.Sprintf("TRUE_%d", i)
	symbol2 := fmt.Sprintf("FINAL_%d", i)
	p.codes = append(p.codes,
		"D=M-D",
		fmt.Sprintf("@%s", symbol1),
	)

	var code string
	switch op {
	case "eq":
		code = "D;JEQ"
	case "gt":
		code = "D;JGT"
	case "lt":
		code = "D;JLT"
	}

	p.codes = append(p.codes,
		code,
		fmt.Sprintf("@%s", SymbolSP),
		"A=M-1",
		"M=0",
		fmt.Sprintf("@%s", symbol2),
		"0;JMP",
		fmt.Sprintf("(%s)", symbol1),
		fmt.Sprintf("@%s", SymbolSP),
		"A=M-1",
		"M=-1",
		fmt.Sprintf("(%s)", symbol2),
	)
	return p
}

func (p *Pipeline) loadAddr(reg string, i int) *Pipeline {

	var base string
	switch reg {
	case "local":
		base = SymbolLCL
	case "argument":
		base = SymbolARG
	case "this":
		base = SymbolTHIS
	case "that":
		base = SymbolTHAT
	}

	p.codes = append(p.codes,
		fmt.Sprintf("@%s", base),
		"D=M",
		fmt.Sprintf("@%d", i),
		"A=D+A",
		"D=M",
	)
	return p
}

func (p *Pipeline) loadTemp(i int) *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%d", TempBase+i),
		"D=M",
	)
	return p
}

func (p *Pipeline) loadConst(i int) *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%d", i),
		"D=A",
	)
	return p
}

func (p *Pipeline) loadStatic(name string, i int) *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%s.%d", name, i),
		"D=M",
	)
	return p
}

func (p *Pipeline) loadPointer(i int) *Pipeline {
	var base string
	switch i {
	case 0:
		base = SymbolTHIS
	case 1:
		base = SymbolTHAT
	}

	p.codes = append(p.codes,
		fmt.Sprintf("@%s", base),
		"D=M",
	)
	return p
}

func (p *Pipeline) loadSP() *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%s", SymbolSP),
		"A=M",
		"D=M",
	)
	return p
}

func (p *Pipeline) storeSP() *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%s", SymbolSP),
		"A=M",
		"M=D",
	)
	return p
}

func (p *Pipeline) computeAddr(reg string, i int) *Pipeline {
	var base string
	switch reg {
	case "local":
		base = SymbolLCL
	case "argument":
		base = SymbolARG
	case "this":
		base = SymbolTHIS
	case "that":
		base = SymbolTHAT
	}

	p.codes = append(p.codes,
		fmt.Sprintf("@%s", base),
		"D=M",
		fmt.Sprintf("@%d", i),
		"D=D+A",
		fmt.Sprintf("@%s", SymbolSP),
		"A=M",
		"M=D",
	)
	return p
}

func (p *Pipeline) storeAddr() *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%s", SymbolSP),
		"A=M+1",
		"A=M",
		"M=D",
	)
	return p
}

func (p *Pipeline) storeStatic(name string, i int) *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%s.%d", name, i),
		"M=D",
	)
	return p
}

func (p *Pipeline) storeTemp(i int) *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%d", TempBase+i),
		"M=D",
	)
	return p
}

func (p *Pipeline) storePointer(i int) *Pipeline {
	var base string
	switch i {
	case 0:
		base = SymbolTHIS
	case 1:
		base = SymbolTHAT
	}

	p.codes = append(p.codes,
		fmt.Sprintf("@%s", base),
		"M=D",
	)
	return p
}

func (p *Pipeline) incrSP() *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%s", SymbolSP),
		"M=M+1",
	)
	return p
}

func (p *Pipeline) decrSP() *Pipeline {
	p.codes = append(p.codes,
		fmt.Sprintf("@%s", SymbolSP),
		"M=M-1",
	)
	return p
}

func (p *Pipeline) label(label string) *Pipeline {
	p.codes = append(p.codes, fmt.Sprintf("(%s)", label))
	return p
}

func (p *Pipeline) Goto(label string) *Pipeline {
	p.codes = append(p.codes,
		"@"+label,
		"0;JMP",
	)
	return p
}

func (p *Pipeline) ifGoto(label string) *Pipeline {
	// @SP
	// AM=M-1
	// D=M
	// @LABEL
	// D;JNE
	p.codes = append(p.codes,
		"@SP",
		"AM=M-1",
		"D=M",
		"@"+label,
		"D;JNE",
	)
	return p
}

func (p *Pipeline) push0(i int) *Pipeline {
	// push constant 0 on LCL
	p.codes = append(p.codes,
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=0",
	)
	// no need to pop because the LCL = SP initially
	return p
}

func (p *Pipeline) function(funcName string, n int) *Pipeline {
	/*
		- push n times (local variable) 0
		// push constant 0
		// pop local i
	*/

	p.codes = append(p.codes,
		fmt.Sprintf("(%s)", funcName),
	)

	for i := 0; i < n; i++ {
		p.push0(i)
	}

	return p
}

func (p *Pipeline) call(hostName, funcName string, argNum, i int) *Pipeline {
	/*
		- push return addr
		- push LCL, ARG, THIS, THAT
		- ARG = SP - 5 - nArgs
		- LCL = SP
		- goto function
		- add return label
	*/

	returnAddr := fmt.Sprintf("%s$ret.%d", hostName, i)

	// -- push return address
	// @returnAddr
	// D=A
	// @SP
	// M=M+1
	// A=M-1
	// M=D
	p.codes = append(p.codes,
		"@"+returnAddr,
		"D=A",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	)

	// -- push LCL
	// @LCL
	// D=M
	// @SP
	// M=M+1
	// A=M-1
	// M=D
	p.codes = append(p.codes,
		"@LCL",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	)

	// -- push ARG
	// @ARG
	// D=M
	// @SP
	// M=M+1
	// A=M-1
	// M=D
	p.codes = append(p.codes,
		"@ARG",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	)

	// -- push THIS
	// @THIS
	// D=M
	// @SP
	// M=M+1
	// A=M-1
	// M=D
	p.codes = append(p.codes,
		"@THIS",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	)

	// -- push THAT
	// @THAT
	// D=M
	// @SP
	// M=M+1
	// A=M-1
	// M=D
	p.codes = append(p.codes,
		"@THAT",
		"D=M",
		"@SP",
		"M=M+1",
		"A=M-1",
		"M=D",
	)

	// -- ARG = SP - 5 - nArgs
	// @SP
	// D=M
	// @5
	// D=D-A
	// @i
	// D=D-A
	// @ARG
	// M=D
	p.codes = append(p.codes,
		"@SP",
		"D=M",
		"@5",
		"D=D-A",
		"@"+strconv.Itoa(argNum),
		"D=D-A",
		"@ARG",
		"M=D",
	)

	// -- LCL = SP
	// @SP
	// D=M
	// @LCL
	// M=D
	p.codes = append(p.codes,
		"@SP",
		"D=M",
		"@LCL",
		"M=D",
	)

	// -- goto function
	// @function
	// 0;JMP
	p.codes = append(p.codes,
		"@"+funcName,
		"0;JMP",
	)

	// -- add return label
	p.codes = append(p.codes,
		fmt.Sprintf("(%s)", returnAddr),
	)

	return p
}

func (p *Pipeline) Return() *Pipeline {
	/*
		- endFrame = LCL
		- returnAddr = *(endFrame - 5)
		- *ARG = pop()
		- SP = ARG + 1
		- THAT = *(endFrame - 1)
		- THIS = *(endFrame - 2)
		- ARG = *(endFrame - 3)
		- LCL = *(endFrame - 4)
		- goto returnAddr
	*/

	// -- endFrame = LCL
	// @LCL
	// D=M
	// @endFrame
	// M=D

	p.codes = append(p.codes,
		"@LCL",
		"D=M",
		"@R13",
		"M=D",
	)

	// -- returnAddr = *(endFrame - 5)
	// @endFrame
	// D=M
	// @5
	// A=D-A
	// D=M
	// @returnAddr
	// M=D
	p.codes = append(p.codes,
		"@R13",
		"D=M",
		"@5",
		"A=D-A",
		"D=M",
		"@R14",
		"M=D",
	)

	// -- *ARG = pop()
	// @SP
	// AM=M-1
	// D=M
	// @ARG
	// A=M
	// M=D
	p.codes = append(p.codes,
		"@SP",
		"AM=M-1",
		"D=M",
		"@ARG",
		"A=M",
		"M=D",
	)

	// -- SP = ARG + 1
	// @ARG
	// D=M+1
	// @SP
	// M=D
	p.codes = append(p.codes,
		"@ARG",
		"D=M+1",
		"@SP",
		"M=D",
	)

	// -- THAT = *(endFrame - 1)
	// @endFrame
	// D=M
	// @1
	// A=D-A
	// D=M
	// @THAT
	// M=D
	p.codes = append(p.codes,
		"@R13",
		"D=M",
		"@1",
		"A=D-A",
		"D=M",
		"@THAT",
		"M=D",
	)

	// -- THIS = *(endFrame - 2)
	// @endFrame
	// D=M
	// @2
	// A=D-A
	// D=M
	// @THIS
	// M=D
	p.codes = append(p.codes,
		"@R13",
		"D=M",
		"@2",
		"A=D-A",
		"D=M",
		"@THIS",
		"M=D",
	)

	// -- ARG = *(endFrame - 3)
	// @endFrame
	// D=M
	// @3
	// A=D-A
	// D=M
	// @ARG
	// M=D
	p.codes = append(p.codes,
		"@R13",
		"D=M",
		"@3",
		"A=D-A",
		"D=M",
		"@ARG",
		"M=D",
	)

	// -- LCL = *(endFrame - 4)
	// @endFrame
	// D=M
	// @4
	// A=D-A
	// D=M
	// @LCL
	// M=D
	p.codes = append(p.codes,
		"@R13",
		"D=M",
		"@4",
		"A=D-A",
		"D=M",
		"@LCL",
		"M=D",
	)

	// -- goto returnAddr
	// @returnAddr
	// A=M
	// 0;JMP
	p.codes = append(p.codes,
		"@R14",
		"A=M",
		"0;JMP",
	)

	return p
}
