package main

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrInvalidOp = fmt.Errorf("invalid operator")
	ErrNumToken  = fmt.Errorf("the number of tokens is wrong")
)

type InstructionType int

const (
	InstructionTypeInvalid InstructionType = iota
	InstructionTypeArithmetic
	InstructionTypeMemory
	InstructionTypeControl
	InstructionTypeFunction
)

type Instruction struct {
	Type     InstructionType
	Cmd      string
	Segment  string
	Offset   int
	Label    string
	HostName string
	FuncName string
	ArgNum   int
}

type Parser struct {
	funcChain []string
}

func NewParser() Parser {
	return Parser{}
}

func (p *Parser) appendFunc(key string) {
	p.funcChain = append(p.funcChain, key)
}

func (p *Parser) popFunc() {
	p.funcChain = p.funcChain[:len(p.funcChain)-1]
}

func (p *Parser) currentFunc() string {
	return p.funcChain[len(p.funcChain)-1]
}

func (p *Parser) hasFunc() bool {
	return len(p.funcChain) > 0
}

func (p *Parser) check(tokens []string) (InstructionType, error) {
	switch tokens[0] {
	case "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not":
		if len(tokens) != 1 {
			return InstructionTypeInvalid, ErrNumToken
		}
		return InstructionTypeArithmetic, nil
	case "push", "pop":
		if len(tokens) != 3 {
			return InstructionTypeInvalid, ErrNumToken
		}
		return InstructionTypeMemory, nil
	case "label", "goto", "if-goto":
		if len(tokens) != 2 {
			return InstructionTypeInvalid, ErrNumToken
		}
		return InstructionTypeControl, nil
	case "function", "call", "return":
		if tokens[0] == "return" && len(tokens) != 1 {
			return InstructionTypeInvalid, ErrNumToken
		}
		if tokens[0] != "return" && len(tokens) != 3 {
			return InstructionTypeInvalid, ErrNumToken
		}
		return InstructionTypeFunction, nil
	default:
		return InstructionTypeInvalid, ErrInvalidOp
	}
}

func (p *Parser) Parse(line string) (Instruction, error) {
	tokens := strings.Split(line, " ")
	if len(tokens) == 0 {
		return Instruction{}, ErrNumToken
	}

	instType, err := p.check(tokens)
	if err != nil {
		fmt.Printf("(%s)\n", line)
		return Instruction{}, ErrInvalidOp
	}

	inst := Instruction{
		Type: instType,
		Cmd:  tokens[0],
	}

	if p.hasFunc() {
		inst.HostName = p.currentFunc()
	}

	if instType == InstructionTypeMemory {
		inst.Segment = tokens[1]
		i, _ := strconv.Atoi(tokens[2])
		inst.Offset = i
	} else if instType == InstructionTypeControl {
		inst.Label = tokens[1]
	} else if instType == InstructionTypeFunction {
		if inst.Cmd != "return" {
			inst.FuncName = tokens[1]
			i, _ := strconv.Atoi(tokens[2])
			inst.ArgNum = i
		}
		if inst.Cmd == "function" && !p.hasFunc() {
			p.appendFunc(inst.FuncName)
		}
		if inst.Cmd == "call" {
			p.appendFunc(inst.FuncName)
		}
		if inst.Cmd == "return" {
			p.popFunc()
		}
	}

	return inst, nil
}
