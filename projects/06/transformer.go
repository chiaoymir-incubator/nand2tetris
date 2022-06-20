package main

import (
	"strconv"
	"strings"
)

var (
	compMap = map[string]string{
		"0":  "101010",
		"1":  "111111",
		"-1": "111010",
		"D":  "001100",
		"A":  "110000", "M": "110000",
		"!D": "001101",
		"!A": "110001", "!M": "110001",
		"-D": "001111",
		"-A": "110011", "-M": "110011",
		"D+1": "011111",
		"A+1": "110111", "M+1": "110111",
		"D-1": "001110",
		"A-1": "110010", "M-1": "110010",
		"D+A": "000010", "D+M": "000010",
		"D-A": "010011", "D-M": "010011",
		"A-D": "000111", "M-D": "000111",
		"D&A": "000000", "D&M": "000000",
		"D|A": "010101", "D|M": "010101",
	}

	destMap = map[string]string{
		"":   "000",
		"M":  "001",
		"D":  "010",
		"DM": "011", "MD": "011",
		"A":  "100",
		"AM": "101", "MA": "101",
		"AD": "110", "DA": "110",
		"ADM": "111", "AMD": "111", "DAM": "111", "DMA": "111", "MAD": "111", "MDA": "111",
	}

	jumpMap = map[string]string{
		"":    "000",
		"JGT": "001",
		"JEQ": "010",
		"JGE": "011",
		"JLT": "100",
		"JNE": "101",
		"JLE": "110",
		"JMP": "111",
	}
)

type Transformer struct {
	Instructions []Instruction
	compMap      map[string]string
	destMap      map[string]string
	jumpMap      map[string]string
}

func NewTransformer(instructions []Instruction) Transformer {
	return Transformer{
		Instructions: instructions,
		compMap:      compMap,
		destMap:      destMap,
		jumpMap:      jumpMap,
	}
}

func (t *Transformer) Transform() []string {
	codes := []string{}

	for _, ins := range t.Instructions {
		if ins.Type == TypeA {
			num, _ := strconv.ParseInt(ins.Tokens[0], 10, 16)
			binNum := strconv.FormatInt(num, 2)
			codes = append(codes, strings.Repeat("0", 16-len(binNum))+binNum)
		} else {
			dest := t.destMap[ins.Tokens[0]]
			comp := t.compMap[ins.Tokens[1]]
			jump := t.jumpMap[ins.Tokens[2]]

			var a string
			if strings.Contains(ins.Tokens[1], "M") {
				a = "1"
			} else {
				a = "0"
			}

			codes = append(codes, "111"+a+comp+dest+jump)
		}
	}

	return codes
}
