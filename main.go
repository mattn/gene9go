package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/text/encoding/japanese"
)

func dictpath() string {
	dir := os.Getenv("HOME")
	if dir == "" && runtime.GOOS == "windows" {
		dir = os.Getenv("APPDATA")
		if dir == "" {
			dir = filepath.Join(os.Getenv("USERPROFILE"), "Application Data", "gene9go")
		}
		dir = filepath.Join(dir, "gene9go")
	} else {
		dir = filepath.Join(dir, ".config", "gene9go")
	}
	return filepath.Join(dir, "gene.txt")
}

func run() int {
	var all bool
	var file string
	flag.StringVar(&file, "f", dictpath(), "path to gene95.txt")
	flag.BoolVar(&all, "a", false, "output all result")
	flag.Parse()

	f, err := os.Open(file)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
		return 1
	}
	defer f.Close()

	word := strings.Join(flag.Args(), " ")
	scanner := bufio.NewScanner(japanese.ShiftJIS.NewDecoder().Reader(f))
	for scanner.Scan() {
		if scanner.Text() == word {
			if !scanner.Scan() {
				break
			}
			text := scanner.Text()
			if !all {
				if words := strings.Split(text, ","); len(words) > 0 {
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