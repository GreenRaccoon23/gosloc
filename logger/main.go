package logger

import (
	"fmt"
	"log"
	"sort"
)

var ()

// Err logs an error and exits
func Err(err error) {

	log.Fatal(err)
}

// // Counts prints the line counts for each file
// func Counts(counts map[string]int) {
//
// 	for fpath, lines := range counts {
// 		fmt.Printf("%d %s\n", lines, fpath)
// 	}
// }

// Counts prints the line counts for each file
func Counts(counts map[string]int) {

	fpaths := []string{}

	for fpath := range counts {
		fpaths = append(fpaths, fpath)
	}

	sort.Strings(fpaths)

	for _, fpath := range fpaths {
		lines := counts[fpath]
		fmt.Printf("%d %s\n", lines, fpath)
	}
}

// Total prints the file total and line total
func Total(fileTotal int, lineTotal int) {

	fmt.Printf("%d files\n", fileTotal)
	fmt.Printf("%d lines\n", lineTotal)
}
