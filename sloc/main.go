package sloc

import (
	"bufio"
	"io"
	"os"
	"sync"

	"github.com/GreenRaccoon23/gosloc/governor"
)

const (
	newLine = '\n'
	// BYTE    = 1.0 << (10 * iota)
	// KILOBYTE
	// MEGABYTE
	// GIGABYTE
	// TERABYTE
)

type counter struct {
	counts map[string]int
	mutex  *sync.Mutex
	g      *governor.Governor
}

func newCounter() *counter {

	counts := make(map[string]int)
	var mutex sync.Mutex

	c := counter{
		counts: counts,
		mutex:  &mutex,
	}

	return &c
}

// Counts returns the number of lines in each file
func Counts(fpaths []string, concurrency int) (map[string]int, error) {

	c := newCounter()
	err := c.readBatch(fpaths, concurrency)
	if err != nil {
		return nil, err
	}

	return c.counts, nil
}

func (c *counter) readBatch(fpaths []string, concurrency int) error {

	size := len(fpaths)
	g := governor.New(size, concurrency)

	for _, fpath := range fpaths {
		g.Accelerate()
		go c.goRead(fpath, g)
	}

	err := g.Regulate()
	if err != nil {
		return err
	}

	return nil
}

func (c *counter) goRead(fpath string, g *governor.Governor) {

	// fpath = filepath.Clean(fpath)
	lines, err := c.read(fpath)
	c.add(fpath, lines)
	g.Decelerate(err)
}

func (c *counter) add(fpath string, lines int) {

	c.mutex.Lock()
	c.counts[fpath] = lines
	c.mutex.Unlock()
}

func (c *counter) read(fpath string) (int, error) {

	f, err := os.Open(fpath)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	defer r.Reset(nil)

	return c.scan(r)
}

func (c *counter) scan(r *bufio.Reader) (int, error) {

	size := bufio.MaxScanTokenSize
	chunk := make([]byte, size)
	total := 0

	for {
		n, err := r.Read(chunk)
		if err != nil && err != io.EOF {
			return 0, err
		}

		if n < size {
			chunk = chunk[:n]
		}

		total += newLines(chunk)

		if err == io.EOF {
			break
		}
	}

	return total, nil
}

func newLines(haystack []byte) int {

	c := 0

	for _, straw := range haystack {
		if straw == newLine {
			c++
		}
	}

	return c
}

// Totals returns the total files and total lines
func Totals(counts map[string]int) (fileTotal int, lineTotal int) {

	for _, lines := range counts {
		fileTotal++
		lineTotal += lines
	}

	return
}
