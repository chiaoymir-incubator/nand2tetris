package main

import (
	"fmt"
)

var (
	ErrNoElem = fmt.Errorf("no such element")
)

const (
	SymbolSP   = "SP"
	SymbolLCL  = "LCL"
	SymbolARG  = "ARG"
	SymbolTHIS = "THIS"
	SymbolTHAT = "THAT"

	TempBase  = 5
	StackBase = 256
)

type Translator struct {
	name  string
	count int
}

func NewTranslator(fileName string) Translator {
	return Translator{
		name:  fileName,
		count: 0,
	}
}

func (ts *Translator) Translate(inst Instruction) []string {
	var codes []string

	switch inst.Type {
	case InstructionTypeArithmetic:
		codes = ts.arithmetic(inst)
	case InstructionTypeMemory:
		codes = ts.memory(inst)
	}

	return codes
}

func (ts *Translator) arithmetic(inst Instruction) []string {
	p := NewPipeline()

	switch inst.Cmd {
	case "add", "sub", "and", "or":
		p.pop2().arithBin(inst.Cmd)

	case "neg", "not":
		p.pop1().arithUni(inst.Cmd)

	case "eq", "gt", "lt":
		p.pop2().arithComp(inst.Cmd, ts.count)
		ts.count++
	}

	return p.codes
}

func (ts *Translator) memory(inst Instruction) []string {
	p := NewPipeline()

	i := inst.Offset
	switch inst.Cmd {
	case "push":
		switch inst.Segment {
		case "local", "argument", "this", "that":
			// addr = *segmentBase+i; *sp = *addr; sp++;
			p.loadAddr(inst.Segment, i)

		case "constant":
			// *sp = i; sp++
			p.loadConst(i)

		case "static":
			// addr = &Name.i; *sp = *addr; sp++;
			p.loadStatic(ts.name, i)

		case "temp":
			// addr = tempBase+i; *sp = *addr; sp++;
			p.loadTemp(i)

		case "pointer":
			// *sp = THIS/THAT; sp++;
			p.loadPointer(i)
		}

		p.storeSP().incrSP()

	case "pop":
		switch inst.Segment {
		case "local", "argument", "this", "that":
			// addr = *segmentBase+i; sp--; *addr = *sp;
			p.computeAddr(inst.Segment, i).decrSP().loadSP().storeAddr()

		case "static":
			// addr = &Name.i; sp--; *addr = *sp;
			p.decrSP().loadSP().storeStatic(ts.name, i)
		case "temp":
			// addr = tempBase+i; sp--; *addr = *sp;
			p.decrSP().loadSP().storeTemp(i)
		case "pointer":
			// sp--; THIS/THAT = *sp;
			p.decrSP().loadSP().storePointer(i)

		case "constant":
			// no-op
		}
	}

	return p.codes
}
