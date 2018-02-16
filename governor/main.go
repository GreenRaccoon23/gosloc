package governor

import (
	"sync"
)

// Governor is a concurrency controlling tool.
// It is passed into goroutines in order to limit how many can run at the same
// time.
// It also pauses control flow until all of the goroutines have finished.
type Governor struct {
	wg          *sync.WaitGroup
	concurrency int
	semaphore   chan bool
	errs        chan error
}

// New creates and initializes a new Governor.
// It will control `size` total goroutines
// by allowing only `concurrency` goroutines to run at the same time.
func New(size int, concurrency int) *Governor {

	if concurrency > size {
		concurrency = size
	}

	var wg sync.WaitGroup
	semaphore := make(chan bool, concurrency)
	errs := make(chan error, size)

	g := Governor{
		wg:          &wg,
		concurrency: concurrency,
		semaphore:   semaphore,
		errs:        errs,
	}

	return &g
}

// Accelerate tells the Governor to control another goroutine.
func (g *Governor) Accelerate() {

	g.wg.Add(1)
	g.semaphore <- true
}

// Decelerate tells the Governor that a goroutine has finished.
func (g *Governor) Decelerate(err error) {

	g.errs <- err
	<-g.semaphore
	g.wg.Done()
}

// Regulate tells the Governor to start watching goroutines and to stop further
// command execution until all of them have finished.
// It returns the first error encountered, if any.
func (g *Governor) Regulate() error {

	g.spin()
	g.coast()
	g.stop()

	return g.condition()
}

func (g *Governor) spin() {

	for govolution := 0; govolution < g.concurrency; govolution++ {
		g.semaphore <- true
	}
}

func (g *Governor) coast() {

	g.wg.Wait()
}

func (g *Governor) stop() {

	close(g.errs)
}

func (g *Governor) condition() error {

	for err := range g.errs {
		if err != nil {
			return err
		}
	}
	return nil
}
