package taskgroup

// A Single manages a single background goroutine. The task is started when the
// value is first created, and the caller can use the Wait method to block
// until it has exited.
type Single[T any] struct {
	valc chan T
	val  T
}

// Wait blocks until the task monitored by s has completed and returns the
// value it reported.
func (s *Single[T]) Wait() T {
	if v, ok := <-s.valc; ok {
		// This is the first call to receive a value:
		// Update val and close the channel (in that order).
		s.val = v
		close(s.valc)
	}
	return s.val
}

// Go runs task in a new goroutine. The caller must call Wait to wait for the
// task to return and collect its value.
func Go[T any](task func() T) *Single[T] {
	// N.B. This is closed by Wait.
	valc := make(chan T, 1)
	go func() { valc <- task() }()

	return &Single[T]{valc: valc}
}

// Run runs task in a new goroutine. The caller must call Wait to wait for the
// task to return. The error reported by Wait is always nil.
func Run(task func()) *Single[error] { return Go(noError(task)) }

// Call starts task in a new goroutine. The caller must call Wait to wait for
// the task to return and collect its result.
func Call[T any](task func() (T, error)) *Single[Result[T]] {
	return Go(func() Result[T] {
		v, err := task()
		return Result[T]{Value: v, Err: err}
	})
}

// A Result is a pair of an arbitrary value and an error.
type Result[T any] struct {
	Value T
	Err   error
}

// Get returns the fields of r as results. It is a convenience method for
// unpacking the results of a Call.
//
// Typical usage:
//
//	s := taskgroup.Call(func() (int, error) { ... })
//	v, err := s.Wait().Get()
func (r Result[T]) Get() (T, error) { return r.Value, r.Err }
