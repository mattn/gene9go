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
	re = regexp.MustCompile(`[;,]`)
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
	var all bool
	var ignorecase bool
	var file string
	flag.StringVar(&file, "f", dictpath(), "path to gene95.txt")
	flag.BoolVar(&all, "a", false, "output all result")
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
		text := scanner.Text()
		if ignorecase {
			text = strings.ToUpper(text)
		}
		if text == word {
			if !scanner.Scan() {
				break
			}
			text = scanner.Text()
			if !all {
				if words := re.Split(text, -1); len(words) > 0 {
					text = words[0]
				}
			}
			fmt.Println(text)
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
