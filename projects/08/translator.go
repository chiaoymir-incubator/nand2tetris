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
	name     string
	count    int
	countMap map[string]int
}

func NewTranslator(fileName string) Translator {
	return Translator{
		name:     fileName,
		count:    0,
		countMap: map[string]int{},
	}
}

// This function initialize the vm implementation segment
func (ts *Translator) Init() []string {
	var codes []string

	codes = append(codes,
		"@261",
		"D=A",
		"@SP",
		"M=D",
		// "@400",
		// "D=A",
		// "@LCL",
		// "M=D",
		// "@500",
		// "D=A",
		// "@ARG",
		// "M=D",
	)

	return codes
}

func (ts *Translator) Translate(inst Instruction) []string {
	var codes []string

	switch inst.Type {
	case InstructionTypeArithmetic:
		codes = ts.arithmetic(inst)
	case InstructionTypeMemory:
		codes = ts.memory(inst)
	case InstructionTypeControl:
		codes = ts.control(inst)
	case InstructionTypeFunction:
		codes = ts.function(inst)
	}

	return codes
}

func (ts *Translator) get(key string) int {
	if v, ok := ts.countMap[key]; !ok {
		return 0
	} else {
		return v
	}
}

func (ts *Translator) add(key string) {
	if _, ok := ts.countMap[key]; !ok {
		ts.countMap[key] = 1
	} else {
		ts.countMap[key] += 1
	}
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

func (ts *Translator) control(inst Instruction) []string {
	p := NewPipeline()

	switch inst.Cmd {
	case "label":
		p.label(inst.Label)
	case "goto":
		p.Goto(inst.Label)
	case "if-goto":
		p.ifGoto(inst.Label)
	}

	return p.codes
}

func (ts *Translator) function(inst Instruction) []string {
	p := NewPipeline()

	switch inst.Cmd {
	case "function":
		p.function(inst.FuncName, inst.ArgNum)
	case "call":
		p.call(inst.HostName, inst.FuncName, inst.ArgNum, ts.get(inst.HostName))
		ts.add(inst.HostName)
	case "return":
		p.Return()
	}

	return p.codes
}
