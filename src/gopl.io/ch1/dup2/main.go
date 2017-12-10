// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 10.
//!+

// Dup2 prints the count and text of lines that appear more than once
// in the input.  It reads from stdin or from a list of named files.
package main

import (
	"bufio"
	"fmt"
	"os"
	"reflect"
)

func main() {
	// counts no of occurences per line, fileNames names of all files where the duplicate occurs
	// inner map of fileNames is a set
	counts, fileNames := make(map[string]int), make(map[string]map[string]struct{})
	files := os.Args[1:]
	if len(files) == 0 {
		countLines(os.Stdin, counts, fileNames)
	} else {
		for _, arg := range files {
			f, err := os.Open(arg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "dup2: %v\n", err)
				continue
			}
			countLines(f, counts, fileNames)
			f.Close()
		}
	}
	for line, n := range counts {
		if n > 1 {
			fmt.Printf("%d\t%s\t%v\n", n, line, reflect.ValueOf(fileNames[line]).MapKeys())
		}
	}
}

func countLines(f *os.File, counts map[string]int, fileNames map[string]map[string]struct{}) {
	input := bufio.NewScanner(f)
	for input.Scan() {
		t := input.Text()
		counts[t]++
		if counts[t] > 1 {
			if fileNames[t] == nil {
				fileNames[t] = make(map[string]struct{})
			}
			fileNames[t][f.Name()] = struct{}{}
		}
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!-
