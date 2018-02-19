package glob

import (
	"os"
	"path/filepath"
)

// globber is a globber
type globber struct {
	rpaths     []string
	inclusions []string
	exclusions []string
	recursive  bool
}

// newGlobber returns a new globber
func newGlobber(rpaths []string, inclusions []string, exclusions []string, recursive bool) globber {

	g := globber{
		rpaths:     rpaths,
		inclusions: inclusions,
		exclusions: exclusions,
		recursive:  recursive,
	}

	return g
}

// Glob runs filepath.Glob, and it does this recursively if requested.
func Glob(rpaths []string, inclusions []string, exclusions []string, recursive bool) ([]string, error) {

	g := newGlobber(rpaths, inclusions, exclusions, recursive)
	matches, err := g.glob()

	return matches, err
}

func (g *globber) glob() ([]string, error) {

	rpaths := g.rpaths
	matches := []string{}

	for _, rpath := range rpaths {

		matches2, err := g.globDynamic(rpath)
		if err != nil {
			return nil, err
		}

		matches = append(matches, matches2...)
	}

	return matches, nil
}

func (g *globber) globDynamic(rpath string) ([]string, error) {

	recursive := g.recursive

	fi, err := os.Lstat(rpath)
	if err != nil {
		return nil, err
	}

	isDir := fi.IsDir()

	if isDir && recursive {
		return g.globRecursive(rpath)
	}

	if isDir && !recursive {
		return g.globThere(rpath)
	}

	return []string{rpath}, nil
}

func (g *globber) globRecursive(rpath string) ([]string, error) {

	matches := []string{}

	dpaths, err := dirsUnder(rpath)
	if err != nil {
		return nil, err
	}

	for _, dpath := range dpaths {
		more, err := g.globThere(dpath)
		if err != nil {
			return nil, err
		}

		matches = append(matches, more...)
	}

	return matches, nil
}

func dirsUnder(rpath string) ([]string, error) {

	dpaths := []string{}

	err := filepath.Walk(rpath, func(fpath string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if fi.IsDir() {
			dpaths = append(dpaths, fpath)
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return dpaths, nil
}

func (g *globber) globThere(dpath string) ([]string, error) {

	cpath, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	defer os.Chdir(cpath)

	err = os.Chdir(dpath)
	if err != nil {
		return nil, err
	}

	matches, err := g.globHere()
	if err != nil {
		return nil, err
	}

	for i := range matches {
		matches[i] = filepath.Join(dpath, matches[i])
	}

	return matches, nil
}

func (g *globber) globHere() ([]string, error) {

	inclusions := g.inclusions
	exclusions := g.exclusions

	if len(inclusions) == 0 {
		inclusions = []string{"*"}
	}

	includes, err := globBatch(inclusions)
	if err != nil {
		return nil, err
	}

	excludes, err := globBatch(exclusions)
	if err != nil {
		return nil, err
	}

	matches := difference(includes, excludes)

	return matches, nil
}

func globBatch(patterns []string) ([]string, error) {

	matches := []string{}

	for _, pattern := range patterns {

		matches2, err := filepath.Glob(pattern)
		if err != nil {
			return nil, err
		}

		matches = append(matches, matches2...)
	}

	return matches, nil
}

func difference(all []string, extras []string) (filtered []string) {

	for _, str := range all {

		if contains(extras, str) {
			continue
		}

		filtered = append(filtered, str)
	}

	return filtered
}

func contains(haystack []string, needle string) bool {

	for _, straw := range haystack {

		if straw == needle {
			return true
		}
	}

	return false
}
