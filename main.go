package main

import (
	"github.com/GreenRaccoon23/gosloc/cmd"
	"github.com/GreenRaccoon23/gosloc/futil"
	"github.com/GreenRaccoon23/gosloc/glob"
	"github.com/GreenRaccoon23/gosloc/logger"
	"github.com/GreenRaccoon23/gosloc/sloc"
)

var ()

func init() {

	cmd.Parse()
}

func main() {

	rpaths := cmd.Rpaths
	inclusions := cmd.Inclusions
	exclusions := cmd.Exclusions
	concurrency := cmd.Concurrency
	total := cmd.Total
	recursive := true

	matches, err := glob.Glob(rpaths, inclusions, exclusions, recursive)
	if err != nil {
		logger.Err(err)
	}

	hardlinks := futil.Hardlinks(matches)

	textpaths, err := futil.TextFiles(hardlinks)
	if err != nil {
		logger.Err(err)
	}

	lineCounts, err := sloc.Counts(textpaths, concurrency)
	if err != nil {
		logger.Err(err)
	}

	fileTotal, lineTotal := sloc.Totals(lineCounts)

	if total {
		logger.Total(fileTotal, lineTotal)
	} else {
		logger.Counts(lineCounts)
	}
}
