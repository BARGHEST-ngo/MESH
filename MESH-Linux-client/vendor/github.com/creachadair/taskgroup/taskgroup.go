// Package taskgroup manages collections of cooperating goroutines.
// It defines a [Group] that handles waiting for goroutine termination and the
// propagation of error values. The caller may provide a callback to filter and
// respond to task errors.
package taskgroup

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
)

// A Task function is the basic unit of work in a [Group]. Errors reported by
// tasks are collected and reported by the group.
type Task func() error

// A Group manages a collection of cooperating goroutines.  Add new tasks to
// the group with [Group.Go] and [Group.Run].  Call [Group.Wait] to wait for
// the tasks to complete. A zero value is ready for use, but must not be copied
// after its first use.
//
// The group collects any errors returned by the tasks in the group. The first
// non-nil error reported by any task (and not otherwise filtered) is returned
// from the Wait method.
type Group struct {
	wg sync.WaitGroup // counter for active goroutines

	// active is nonzero when the group is "active", meaning there has been at
	// least one call to Go since the group was created or the last Wait.
	//
	// Together active and μ work as a kind of resettable sync.Once; the fast
	// path reads active and only acquires μ if it discovers setup is needed.
	active atomic.Uint32

	μ       sync.Mutex // guards the fields below
	err     error      // error returned from Wait
	onError errorFunc  // called each time a task returns non-nil
}

// activate resets the state of the group and marks it as active.  This is
// triggered by adding a goroutine to an empty group.
func (g *Group) activate() {
	g.μ.Lock()
	defer g.μ.Unlock()
	if g.active.Load() == 0 { // still inactive
		g.err = nil
		g.active.Store(1)
	}
}

// New constructs a new empty group with the specified error filter.
// See [Group.OnError] for a description of how errors are filtered.
// If ef == nil, no filtering is performed.
func New(ef any) *Group { return new(Group).OnError(ef) }

// OnError sets the error filter for g. If ef == nil, the error filter is
// removed and errors are no longer filtered. Otherwise, each non-nil error
// reported by a task running in g is passed to ef.
//
// The concrete type of ef must be a function with one of the following
// signature schemes, or OnError will panic.
//
// If ef is:
//
//	func()
//
// then ef is called once per reported error, and the error is not modified.
//
// If ef is:
//
//	func(error)
//
// then ef is called with each reported error, and the error is not modified.
//
// If ef is:
//
//	func(error) error
//
// then ef is called with each reported error, and its result replaces the
// reported value. This permits ef to suppress or replace the error value
// selectively.
//
// Calls to ef are synchronized so that it is safe for ef to manipulate local
// data structures without additional locking. It is safe to call OnError while
// tasks are active in g.
func (g *Group) OnError(ef any) *Group {
	filter := adaptErrorFunc(ef)
	g.μ.Lock()
	defer g.μ.Unlock()
	g.onError = filter
	return g
}

// Go runs task in a new goroutine in g.
func (g *Group) Go(task Task) {
	g.wg.Add(1)
	if g.active.Load() == 0 {
		g.activate()
	}
	go func() {
		defer g.wg.Done()
		if err := task(); err != nil {
			g.handleError(err)
		}
	}()
}

// Run runs task in a new goroutine in g.
// The resulting task reports a nil error.
func (g *Group) Run(task func()) { g.Go(noError(task)) }

func (g *Group) handleError(err error) {
	g.μ.Lock()
	defer g.μ.Unlock()
	e := g.onError.filter(err)
	if e != nil && g.err == nil {
		g.err = e // capture the first unfiltered error always
	}
}

// Wait blocks until all the goroutines currently active in the group have
// returned, and all reported errors have been delivered to the callback.  It
// returns the first non-nil error reported by any of the goroutines in the
// group and not filtered by an OnError callback.
//
// As with [sync.WaitGroup], new tasks can be added to g during a call to Wait
// only if the group contains at least one active task when Wait is called and
// continuously thereafter until the last concurrent call to g.Go returns.
//
// Wait may be called from at most one goroutine at a time.  After Wait has
// returned, the group is ready for reuse.
func (g *Group) Wait() error {
	g.wg.Wait()
	g.μ.Lock()
	defer g.μ.Unlock()

	// If the group is still active, deactivate it now.
	if g.active.Load() != 0 {
		g.active.Store(0)
	}
	return g.err
}

// An errorFunc is called by a group each time a task reports an error.  Its
// return value replaces the reported error, so the errorFunc can filter or
// suppress errors by modifying or discarding the input error.
type errorFunc func(error) error

func (ef errorFunc) filter(err error) error {
	if ef == nil {
		return err
	}
	return ef(err)
}

var (
	triggerType = reflect.TypeOf(func() {})
	listenType  = reflect.TypeOf(func(error) {})
	filterType  = reflect.TypeOf(func(error) error { return nil })
)

func adaptErrorFunc(ef any) errorFunc {
	v := reflect.ValueOf(ef)
	if !v.IsValid() {
		// OK, ef == nil, nothing to do
		return nil
	} else if t := v.Type(); t.ConvertibleTo(triggerType) {
		f := v.Convert(triggerType).Interface().(func())
		return func(err error) error { f(); return err }
	} else if t.ConvertibleTo(listenType) {
		f := v.Convert(listenType).Interface().(func(error))
		return func(err error) error { f(err); return err }
	} else if t.ConvertibleTo(filterType) {
		return errorFunc(v.Convert(filterType).Interface().(func(error) error))
	} else {
		panic(fmt.Sprintf("unsupported filter type %T", ef))
	}
}

// NoError adapts f to a Task that executes f and reports a nil error.
//
// Deprecated: Use [Group.Run] or [taskgroup.Run] instead.
func NoError(f func()) Task { return func() error { f(); return nil } }

func noError(f func()) Task { return func() error { f(); return nil } }

// Limit returns g and a [StartFunc] that starts each task passed to it in g,
// allowing no more than n tasks to be active concurrently. If n ≤ 0, no limit
// is enforced.
//
// The limiting mechanism is optional, and the underlying group is not
// restricted. A call to the start function will block until a slot is
// available, but calling g.Go directly will add a task unconditionally and
// will not take up a limiter slot.
//
// This is a shorthand for constructing a [Throttle] with capacity n and
// calling its Limit method.  If n ≤ 0, the start function is equivalent to
// g.Go, which enforces no limit. To share a throttle among multiple groups,
// construct the throttle separately.
func (g *Group) Limit(n int) (*Group, StartFunc) { t := NewThrottle(n); return g, t.Limit(g) }
