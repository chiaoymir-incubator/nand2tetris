package main

import (
	"log"
	"strings"
)

type InstructionType int

const (
	TypeA InstructionType = iota
	TypeC
)

type Instruction struct {
	Type   InstructionType
	Tokens []string
}

type Parser struct {
}

func NewParser() Parser {
	return Parser{}
}

func (p *Parser) Parse(lines []string) []Instruction {
	instructions := []Instruction{}
	r := strings.NewReplacer("@", "")
	for i, line := range lines {
		var ins Instruction
		if strings.Contains(line, "@") {
			ins = Instruction{Type: TypeA}
			num := r.Replace(line)
			ins.Tokens = append(ins.Tokens, num)
		} else {
			ins = Instruction{Type: TypeC}

			if strings.Contains(line, "=") {
				tokens := strings.Split(line, "=")
				if len(tokens) != 2 {
					log.Fatalf("incorrect command: %d", i)
				}
				ins.Tokens = append(ins.Tokens, tokens[0])
				line = tokens[1]
			} else {
				ins.Tokens = append(ins.Tokens, "")
			}

			if strings.Contains(line, ";") {
				tokens := strings.Split(line, ";")
				if len(tokens) != 2 {
					log.Fatalf("incorrect command: %d", i)
				}
				ins.Tokens = append(ins.Tokens, tokens...)
			} else {
				ins.Tokens = append(ins.Tokens, line, "")
			}
		}
		instructions = append(instructions, ins)
	}

	return instructions
}
