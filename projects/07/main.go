package main

import (
	"bufio"
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	SourceExt string = ".vm"
	TargetExt string = ".asm"
)

var fileName = flag.String("f", "", "the vm file")

func process(lines []string, w *bufio.Writer, parser Parser, translator Translator) {
	for _, line := range lines {
		// Remove space and comments
		line, ok := sanitize(line)
		if !ok {
			continue
		}

		// Parse line to tokens
		inst, err := parser.Parse(line)
		if err != nil {
			log.Fatalln(err)
		}

		// Translation
		codes := translator.Translate(inst)

		w.WriteString("// " + line + "\n")
		for _, code := range codes {
			w.WriteString(code + "\n")
		}
	}

	if err := w.Flush(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	flag.Parse()

	base := strings.TrimSuffix(*fileName, filepath.Ext(*fileName))
	if base == "" {
		log.Printf("The program needs an argument: -f <file>")
		os.Exit(1)
	}

	parser := NewParser()
	translator := NewTranslator(filepath.Base(base))

	// Read from file
	file, err := os.Open(base + SourceExt)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	f, err := os.Create(base + TargetExt)
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}

	c := 0
	scanner := bufio.NewScanner(file)
	lines := make([]string, 100)
	for scanner.Scan() {
		// Read 100 lines each time
		lines = append(lines, scanner.Text())
		c++
		if c%100 == 0 {
			process(lines, w, parser, translator)
			c = 0
			lines = make([]string, 100)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// process the remaining lines
	process(lines, w, parser, translator)
}
