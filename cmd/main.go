package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/GreenRaccoon23/gosloc/futil"
	"github.com/GreenRaccoon23/gosloc/logger"
	"github.com/spf13/pflag"
)

const (
	newLine = '\n'
)

var (
	// inclusionsBunch description under parsArgs
	inclusionsBunch string
	// exclusionsBunch description under parsArgs
	exclusionsBunch string
	// Concurrency description under parsArgs
	Concurrency int
	// Total description under parsArgs
	Total bool
	// Rpaths description under parsArgs
	Rpaths []string

	// Inclusions is inclusionsBunch split by ','
	Inclusions []string
	// Exclusions is inclusionsBunch split by ','
	Exclusions []string
)

// Parse parses command line arguments
func Parse() {

	parsArgs()
	parseStdin()
	// debug()
	// os.Exit(0)
	validate()
	setInclusions()
	setExclusions()
}

func parsArgs() {

	pflag.StringVarP(&inclusionsBunch, "include", "i", "", "File patterns to include, separated by commas")
	pflag.StringVarP(&exclusionsBunch, "exclude", "x", "", "File patterns to exclude, separated by commas")
	pflag.IntVarP(&Concurrency, "concurrency", "c", 1, "Max number of files to read simultaneously")
	pflag.BoolVarP(&Total, "total", "t", false, "Show a grand total, not the total for each file")
	pflag.Usage = usage
	pflag.CommandLine.SortFlags = false
	pflag.Parse()
	Rpaths = pflag.Args()
}

func parseStdin() {

	stdin := os.Stdin

	if !futil.IsPipe(stdin) {
		return
	}

	rpaths, err := futil.ReadLines(stdin)
	if err != nil {
		logger.Err(err)
	}

	Rpaths = append(Rpaths, rpaths...)
}

// func debug() {
//
// 	fmt.Printf("inclusionsBunch: %v\n", inclusionsBunch)
// 	fmt.Printf("exclusionsBunch: %v\n", exclusionsBunch)
// 	fmt.Printf("Concurrency: %v\n", Concurrency)
// 	fmt.Printf("Total: %v\n", Total)
// 	fmt.Printf("Rpaths: %v\n", Rpaths)
// }

// usage overrides pflag.Usage
func usage() {
	fmt.Fprintf(os.Stderr, "gosloc <options> <path>...\n")
	pflag.PrintDefaults()
	fmt.Fprintf(os.Stderr,
		`
WARNING: Setting concurrency too high will cause the program to crash.
`,
	)
}

func validate() {

	if len(Rpaths) == 0 {
		complain("No paths specified")
	}

	if Concurrency <= 0 {
		complain("-c (concurrency) must be above 0")
	}
}

func complain(complaint string) {

	fmt.Fprintf(os.Stderr, "%v\n\n", complaint)
	usage()
	fmt.Fprintf(os.Stderr, "\n%v\n", complaint)
	os.Exit(2)
}

func setInclusions() {

	if inclusionsBunch == "" {
		return
	}

	Inclusions = strings.Split(inclusionsBunch, ",")
}

func setExclusions() {

	if exclusionsBunch == "" {
		return
	}

	Exclusions = strings.Split(exclusionsBunch, ",")
}
