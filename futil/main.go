package futil

import (
	"bufio"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type sliceHash map[string]bool

var (
	textFileExts = newSliceHash(
		"asm",
		"brs",
		"c",
		"cc",
		"clj",
		"cljs",
		"coffee",
		"cpp",
		"cr",
		"cs",
		"csh",
		"css",
		"cxx",
		"erl",
		"f",
		"go",
		"groovy",
		"gs",
		"h",
		"handlebars", "hbs",
		"hpp",
		"hr",
		"hs",
		"html", "htm",
		"hx",
		"hxx",
		"hy",
		"iced",
		"ino",
		"jade",
		"java",
		"jl",
		"js",
		"jsx",
		"kt",
		"kts",
		"less",
		"lua",
		"ls",
		"m",
		"md", "markdown",
		"mm",
		"ml",
		"mli",
		"mochi",
		"monkey",
		"mustache",
		"nix",
		"nim",
		"nut",
		"php", "php5",
		"pl",
		"py",
		"r",
		"rb",
		"rkt",
		"rs",
		"sass",
		"scala",
		"scss",
		"sh",
		"styl",
		"svg",
		"sql",
		"swift",
		"ts",
		"tsx",
		"txt", "text",
		"vb",
		"xhtml",
		"xml",
		"yaml",
	)
)

func newSliceHash(strs ...string) *sliceHash {

	sh := sliceHash{}

	for _, s := range strs {
		sh[s] = true
	}

	return &sh
}

func (sh *sliceHash) has(str string) bool {

	_, has := (*sh)[str]

	return has
}

// IsPipe returns whether a file is a pipe
func IsPipe(f *os.File) bool {

	fi, err := f.Stat()
	if err != nil {
		return false
	}

	mode := fi.Mode()

	if mode&os.ModeNamedPipe != 0 {
		return true
	}

	return false
}

// ReadLines reads a file and returns its lines as a slice
func ReadLines(f *os.File) ([]string, error) {

	lines := []string{}

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		lines = append(lines, line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

// AnyHardlinks checks whether any of the fpaths points to a hardlink
func AnyHardlinks(fpaths []string) bool {

	for _, fpath := range fpaths {

		if isHardlink(fpath) {
			return true
		}
	}

	return false
}

// Hardlinks returns only the hardlinks in fpaths
// (i.e., non-directories and non-symlinks)
func Hardlinks(fpaths []string) []string {

	filtered := []string{}

	for _, fpath := range fpaths {

		if !isHardlink(fpath) {
			continue
		}

		filtered = append(filtered, fpath)
	}

	return filtered
}

func isHardlink(fpath string) bool {

	fi, err := os.Lstat(fpath)
	if err != nil {
		return false
	}

	mode := fi.Mode()

	if !mode.IsRegular() {
		return false
	}

	if mode&os.ModeTemporary != 0 {
		return false
	}

	if mode&os.ModeCharDevice != 0 {
		return false
	}

	return true
}

// TextFiles returns only the text files in fpaths
func TextFiles(fpaths []string) ([]string, error) {

	filtered := []string{}

	for _, fpath := range fpaths {

		is, err := isTextFile(fpath)
		if err != nil {
			return nil, err
		}

		if !is {
			continue
		}

		filtered = append(filtered, fpath)
	}

	return filtered, nil
}

func isTextFile(fpath string) (bool, error) {

	if isTextFileByExt(fpath) {
		return true, nil
	}

	return isTextFileByHeader(fpath)
}

func isTextFileByExt(fpath string) bool {

	ext := filepath.Ext(fpath)
	ext = strings.TrimPrefix(ext, ".")

	if textFileExts.has(ext) {
		return true
	}

	return false
}

func isTextFileByHeader(fpath string) (bool, error) {

	essence, err := mimeTypeByHeader(fpath) // "text/plain"
	if err != nil {
		return false, err
	}

	mainType := strBefore(essence, "/") // "text"

	if mainType == "text" {
		return true, nil
	}

	// switch mainType {
	// case "text":
	// 	return true, nil
	//
	// case "application":
	// 	break
	//
	// case "image":
	// 	fallthrough
	// case "audio":
	// 	fallthrough
	// case "font":
	// 	fallthrough
	// case "video":
	// 	fallthrough
	// case "model":
	// 	fallthrough
	// case "chemical":
	// 	fallthrough
	// case "x-conference":
	// 	fallthrough
	// case "message":
	// 	return false, nil
	// }

	return false, nil
}

func mimeTypeByHeader(fpath string) (string, error) {

	f, err := os.Open(fpath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	chunk := make([]byte, 512)

	_, err = f.Read(chunk)
	if err == io.EOF {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	mimeType := http.DetectContentType(chunk) // "text/plain; charset=utf-8"
	essence := strBefore(mimeType, ";")       // "text/plain"

	return essence, nil
}

func strBefore(haystack string, needle string) string {

	i := strings.Index(haystack, needle)
	if i == -1 {
		return haystack
	}

	return haystack[:i]
}

// func slcContains(haystack []string, needle string) bool {
//
// 	for _, straw := range haystack {
// 		if straw == needle {
// 			return true
// 		}
// 	}
//
// 	return false
// }
