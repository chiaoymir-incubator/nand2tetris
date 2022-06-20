package main

import "fmt"

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
