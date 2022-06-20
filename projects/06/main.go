package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var fileName = flag.String("f", "", "the assembly file")

func main() {
	flag.Parse()

	base := strings.TrimSuffix(*fileName, filepath.Ext(*fileName))
	if base == "" {
		log.Printf("The program needs an argument: -f <file>")
		os.Exit(1)
	}

	// Read from file
	file, err := os.Open(base + ".asm")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	lines := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line, ok := sanitize(scanner.Text())
		if ok {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Translation
	st := NewSymbolTable()
	lines = st.Replace(lines)

	parser := NewParser()
	instructions := parser.Parse(lines)

	transformer := NewTransformer(instructions)
	codes := transformer.Transform()

	// Write to file
	f, err := os.Create(base + ".hack")
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}

	for _, code := range codes {
		w.WriteString(code + "\n")
	}

	if err := w.Flush(); err != nil {
		log.Fatal(err)
	}
}
