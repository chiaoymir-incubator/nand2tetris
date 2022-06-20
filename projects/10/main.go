package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	SourceExt     string = ".jack"
	TargetExt     string = ".xml"
	SourcePattern string = "*.jack"
)

var fileName = flag.String("f", "", "the jack file")

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

func compile(fName string) {
	// create target file
	name := strings.TrimSuffix(fName, SourceExt)
	f, err := os.Create(name + "_c" + TargetExt)
	if err != nil {
		log.Printf("Error create file: %s\n", name)
		return
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	// open source file
	f, err = os.Open(fName)
	if err != nil {
		log.Printf("Error open file: %s\n", fName)
		return
	}
	defer f.Close()

	sanitizer := Sanitizer{}
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		// Remove space and comments
		sanitizer.s = scanner.Text()
		line, ignore := sanitizer.Sanitize()
		if ignore {
			continue
		}
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
		return
	}

	sanitizer.s = strings.Join(lines, "")
	str, _ := sanitizer.removeComments()
	tokenizer := Tokenizer{str: str}

	/* This is for tokenizer testing */
	// tokens := []Token{}
	// tokenizer.Advance()
	// for tokenizer.HasMoreTokens() {
	// 	token := tokenizer.GetToken()
	// 	tokens = append(tokens, token)
	// 	tokenizer.Advance()
	// }
	// fmt.Println()
	// var strs []string
	// for _, token := range tokens {
	// 	strs = append(strs, token.content)
	// }
	// fmt.Println(strings.Join(strs, " "))

	// w.WriteString("<tokens>\n")
	// for _, token := range tokens {
	// 	w.WriteString(fmt.Sprintf("<%s> %s </%s>\n", token.tokenType, token.content, token.tokenType))
	// }
	// w.WriteString("</tokens>\n")

	ce := NewCompilationEngine(tokenizer, w)
	if err := ce.Compile(); err != nil {
		log.Printf("Error parse file: %s. Error: %s.\n", fName, err.Error())

		// for debugging: flust content to see parsing process
		if err := w.Flush(); err != nil {
			log.Println(err)
		}
		return
	}

	if err := w.Flush(); err != nil {
		log.Println(err)
	}
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
	// var targetDir, targetName string
	if fileInfo.IsDir() {
		fileNames, err = WalkMatch(*fileName, SourcePattern)
		if err != nil {
			log.Fatalln(err)
		}
	} else {
		fileNames = append(fileNames, *fileName)
	}

	// intentionally close the file or directory
	file.Close()

	var wg sync.WaitGroup
	wg.Add(len(fileNames))
	// merge all files into a single file
	for _, fName := range fileNames {
		go func(fName string) {
			start := time.Now()
			compile(fName)
			lapse := float64(time.Since(start).Nanoseconds()) / 1000000.0
			fmt.Printf("Processed file: %s.\tTime lapse: %f ms.\n", fName, lapse)
			wg.Done()
		}(fName)
	}
	wg.Wait()
}
