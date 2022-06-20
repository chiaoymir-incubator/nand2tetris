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
)

type Instruction struct {
	Type    InstructionType
	Cmd     string
	Segment string
	Offset  int
}

type Parser struct {
}

func NewParser() Parser {
	return Parser{}
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
		return Instruction{}, ErrInvalidOp
	}

	inst := Instruction{
		Type: instType,
		Cmd:  tokens[0],
	}

	if instType == InstructionTypeMemory {
		inst.Segment = tokens[1]
		i, _ := strconv.Atoi(tokens[2])
		inst.Offset = i
	}

	return inst, nil
}
