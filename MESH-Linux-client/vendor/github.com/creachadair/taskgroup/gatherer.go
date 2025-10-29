package taskgroup

import "sync"

// A Gatherer manages a group of [Task] functions that report values, and
// gathers the values they return.
type Gatherer[T any] struct {
	run func(Task) // start the task in a goroutine

	μ      sync.Mutex
	gather func(T) // handle values reported by tasks
}

func (g *Gatherer[T]) report(v T) {
	g.μ.Lock()
	defer g.μ.Unlock()
	g.gather(v)
}

// Gather creates a new empty gatherer that uses run to execute tasks returning
// values of type T.
//
// If gather != nil, values reported by successful tasks are passed to the
// function, otherwise such values are discarded. Calls to gather are
// synchronized to a single goroutine.
//
// If run == nil, Gather will panic.
func Gather[T any](run func(Task), gather func(T)) *Gatherer[T] {
	if run == nil {
		panic("run function is nil")
	}
	if gather == nil {
		gather = func(T) {}
	}
	return &Gatherer[T]{run: run, gather: gather}
}

// Call runs f in g. If f reports an error, the error is propagated to the
// runner; otherwise the non-error value reported by f is gathered.
func (g *Gatherer[T]) Call(f func() (T, error)) {
	g.run(func() error {
		v, err := f()
		if err == nil {
			g.report(v)
		}
		return err
	})
}

// Run runs f in g, and gathers the value it reports.
func (g *Gatherer[T]) Run(f func() T) {
	g.run(func() error { g.report(f()); return nil })
}

// Report runs f in g. Any values passed to report are gathered.  If f reports
// an error, that error is propagated to the runner.  Any values sent before f
// returns are still gathered, even if f reports an error.
func (g *Gatherer[T]) Report(f func(report func(T)) error) {
	g.run(func() error { return f(g.report) })
}
