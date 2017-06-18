package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

var (
	resplit = regexp.MustCompile(`[;,]`)
	renorm  = regexp.MustCompile(`^[0-9]+\.`)
)

func dictpath() string {
	dir := ""
	if runtime.GOOS == "windows" {
		dir = filepath.Join(os.Getenv("APPDATA"), "gene9go")
	} else {
		dir = filepath.Join(os.Getenv("HOME"), ".config", "gene9go")
	}
	return filepath.Join(dir, "gene.txt")
}

func run() int {
	var only bool
	var ignorecase bool
	var file string
	flag.StringVar(&file, "f", dictpath(), "path to gene95.txt")
	flag.BoolVar(&only, "o", false, "show first candidate")
	flag.BoolVar(&ignorecase, "i", false, "ignore case")
	flag.Parse()

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	defer f.Close()

	word := strings.Join(flag.Args(), " ")
	if ignorecase {
		word = strings.ToUpper(word)
	}
	scanner := bufio.NewScanner(japanese.ShiftJIS.NewDecoder().Reader(f))
	for scanner.Scan() {
		first := scanner.Text()
		if ignorecase {
			first = strings.ToUpper(first)
		}
		if !scanner.Scan() {
			break
		}
		second := scanner.Text()

		if first == word {
			if only {
				if words := resplit.Split(second, -1); len(words) > 0 {
					second = words[0]
				}
				second = renorm.ReplaceAllString(second, "")
			}
			fmt.Println(second)
			return 0
		}
	}
	if scanner.Err() != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	return 2
}

func main() {
	os.Exit(run())
}
