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

// https://stackoverflow.com/questions/55300117/how-do-i-find-all-files-that-have-a-certain-extension-in-go-regardless-of-depth
func WalkMatch(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func getEntry(fileNames []string, fileName string) string {
	for _, v := range fileNames {
		if filepath.Base(v) == fileName {
			return v
		}
	}
	return ""
}

func removeEntry(fileNames []string, fileName string) []string {
	newFileNames := []string{}
	for _, v := range fileNames {
		if v != fileName {
			newFileNames = append(newFileNames, v)
		}
	}
	return newFileNames
}

func sortEntries(fileNames []string) []string {
	newFileNames := []string{}
	if fName := getEntry(fileNames, "Sys.vm"); fName != "" {
		newFileNames = append(newFileNames, fName)
		removedEntries := removeEntry(fileNames, fName)
		newFileNames = append(newFileNames, removedEntries...)
	} else {
		newFileNames = fileNames
	}

	return newFileNames
}

func main() {
	flag.Parse()

	// Read from file
	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatalln(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatalln(err)
	}

	var fileNames []string
	var targetDir, targetName string
	if fileInfo.IsDir() {
		fileNames, err = WalkMatch(*fileName, "*.vm")
		if err != nil {
			log.Fatalln(err)
		}
		targetDir = *fileName
		targetName = filepath.Base(*fileName)
	} else {
		fileNames = append(fileNames, *fileName)
		targetDir = filepath.Dir(*fileName)
		targetName = filepath.Base(targetDir)
	}

	fileNames = sortEntries(fileNames)

	// intentionally
	file.Close()

	f, err := os.Create(filepath.Join(targetDir, targetName+TargetExt))
	if err != nil {
		panic(err)
	}

	w := bufio.NewWriter(f)
	if err != nil {
		panic(err)
	}

	parser := NewParser()
	translator := NewTranslator(targetName)

	// Note:
	// This bootstrap code should only be used in the following two test
	// 1. FibonacciElement
	// 2. StaticsTest
	// You can comment out these lines in other tests
	initCodes := translator.Init()
	w.WriteString("// initialization \n")
	for _, code := range initCodes {
		w.WriteString(code + "\n")
	}

	// merge all files into a single file
	for _, fName := range fileNames {
		name := strings.TrimSuffix(fName, filepath.Ext(fName))
		translator.name = filepath.Base(name)

		f, err := os.Open(fName)
		if err != nil {
			log.Fatalln(err)
		}

		scanner := bufio.NewScanner(f)
		for scanner.Scan() {
			// Remove space and comments
			line, ok := sanitize(scanner.Text())
			if !ok {
				continue
			}

			// Parse a line to an instruction
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

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	}

	if err := w.Flush(); err != nil {
		log.Fatal(err)
	}
}
