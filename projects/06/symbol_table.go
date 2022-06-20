package main

import (
	"fmt"
	"strconv"
	"strings"
)

var (
	ErrNoElem = fmt.Errorf("no such element")
)

var (
	registerBase = 16
)

type SymbolTable struct {
	m          map[string]int
	labelCount int
	varCount   int
}

func NewSymbolTable() SymbolTable {
	m := map[string]int{
		"SCREEN": 16384,
		"KBD":    24576,
		"SP":     0,
		"LCL":    1,
		"ARG":    2,
		"THIS":   3,
		"THAT":   4,
	}

	for i := 0; i < registerBase; i++ {
		m[fmt.Sprintf("R%d", i)] = i
	}

	return SymbolTable{
		m:          m,
		labelCount: 0,
		varCount:   registerBase,
	}
}

func (st *SymbolTable) Replace(lines []string) []string {
	lines = st.Scan(lines)
	newLines := []string{}
	r := strings.NewReplacer("@", "")

	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "@"):
			symbol := r.Replace(line)

			// not a number
			if !isInt(symbol) {
				val, err := st.Get(symbol)
				if err != nil {
					fmt.Println(symbol)
					panic(err)
				}
				newLines = append(newLines, fmt.Sprintf("@%d", val))
				continue
			}
		}
		newLines = append(newLines, line)
	}

	return newLines
}

func (st *SymbolTable) Scan(lines []string) []string {
	return st.ScanVar(st.ScanLabel(lines))
}

func (st *SymbolTable) ScanLabel(lines []string) []string {
	st.labelCount = 0
	newLines := []string{}
	r := strings.NewReplacer("(", "", ")", "", "@", "")
	for i, line := range lines {
		switch {
		case strings.HasPrefix(line, "("):
			symbol := r.Replace(line)
			if !st.Contains(symbol) {
				st.Set(symbol, i-st.labelCount)
				st.labelCount++
			}
			continue
		}
		newLines = append(newLines, line)
	}

	return newLines
}

func (st *SymbolTable) ScanVar(lines []string) []string {
	st.varCount = registerBase
	newLines := []string{}
	r := strings.NewReplacer("@", "")
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "@"):
			symbol := r.Replace(line)

			// not a number
			if !isInt(symbol) {
				if !st.Contains(symbol) {
					st.Set(symbol, st.varCount)
					st.varCount++
				}
			}

		}
		newLines = append(newLines, line)
	}

	return newLines
}

func (st *SymbolTable) Contains(symbol string) bool {
	_, ok := st.m[symbol]
	return ok
}

func (st *SymbolTable) Get(symbol string) (int, error) {
	if !st.Contains(symbol) {
		return 0, ErrNoElem
	}

	return st.m[symbol], nil
}

func (st *SymbolTable) Set(symbol string, val int) {
	st.m[symbol] = val
}

func (st *SymbolTable) Print() {
	for k, v := range st.m {
		fmt.Printf("%s: %d\n", k, v)
	}
}

func isInt(symbol string) bool {
	_, err := strconv.ParseInt(symbol, 10, 64)
	return err == nil
}
