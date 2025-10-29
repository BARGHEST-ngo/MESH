package taskgroup

import "sync/atomic"

// A Throttle rate-limits the number of concurrent goroutines that can execute
// in parallel to some fixed number.  A zero Throttle is ready for use, but
// imposes no limit on parallel execution.
type Throttle struct {
	adm chan struct{}
}

// NewThrottle constructs a [Throttle] with a capacity of n goroutines.
// If n ≤ 0, the resulting Throttle imposes no limit.
func NewThrottle(n int) Throttle {
	if n <= 0 {
		return Throttle{}
	}
	return Throttle{adm: make(chan struct{}, n)}
}

// enter blocks until a slot is available in t, then returns a func that the
// caller must execute to return the slot when it is no longer in use.
func (t Throttle) enter() func() {
	if t.adm == nil {
		return func() {}
	}
	t.adm <- struct{}{}
	var done atomic.Bool
	return func() {
		if done.CompareAndSwap(false, true) {
			<-t.adm
		}
	}
}

// Limit returns a function that starts each [Task] passed to it in g,
// respecting the rate limit imposed by t. Each call to Limit yields a fresh
// start function, and all the functions returned share the capacity of t.
func (t Throttle) Limit(g *Group) StartFunc {
	return func(task Task) {
		release := t.enter()
		g.Go(func() error {
			defer release()
			return task()
		})
	}
}

// A StartFunc executes each [Task] passed to it in a [Group].
type StartFunc func(Task)

// Go is a legibility shorthand for calling s with task.
func (s StartFunc) Go(task Task) { s(task) }

// Run is a legibility shorthand for calling s with a task that runs f and
// reports a nil error.
func (s StartFunc) Run(f func()) { s(noError(f)) }
