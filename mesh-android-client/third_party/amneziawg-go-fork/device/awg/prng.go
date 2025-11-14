package awg

import (
	crand "crypto/rand"
	"fmt"
	v2 "math/rand/v2"

	"golang.org/x/exp/constraints"
)

type RandomNumberGenerator[T constraints.Integer] interface {
	RandomSizeInRange(min, max T) T
	Get() uint64
	ReadSize(size int) []byte
}

type PRNG[T constraints.Integer] struct {
	cha8Rand *v2.ChaCha8
}

func NewPRNG[T constraints.Integer]() *PRNG[T] {
	buf := make([]byte, 32)
	n, err := crand.Read(buf)
	if err != nil || n != 32 {
		panic(fmt.Sprintf("NewPRNG: failed to read random bytes: err=%v, n=%d", err, n))
	}

	cha8 := v2.NewChaCha8([32]byte(buf))
	if cha8 == nil {
		panic("NewPRNG: NewChaCha8 returned nil")
	}

	return &PRNG[T]{
		cha8Rand: cha8,
	}
}

func (p *PRNG[T]) RandomSizeInRange(min, max T) T {
	// DEBUG: Check if cha8Rand is initialized
	if p == nil || p.cha8Rand == nil {
		panic("PRNG.RandomSizeInRange: PRNG or cha8Rand is nil - PRNG was not initialized with NewPRNG()")
	}

	if min > max {
		panic("min must be less than max")
	}

	if min == max {
		return min
	}

	return T(p.Get()%uint64(max-min)) + min
}

func (p *PRNG[T]) Get() uint64 {
	// DEBUG: Check if cha8Rand is initialized
	if p == nil || p.cha8Rand == nil {
		panic("PRNG.Get: PRNG or cha8Rand is nil - PRNG was not initialized with NewPRNG()")
	}

	return p.cha8Rand.Uint64()
}

func (p *PRNG[T]) ReadSize(size int) []byte {
	// DEBUG: Check if cha8Rand is initialized
	if p == nil || p.cha8Rand == nil {
		panic("PRNG.ReadSize: PRNG or cha8Rand is nil - PRNG was not initialized with NewPRNG()")
	}

	// TODO: use a memory pool to allocate
	buf := make([]byte, size)
	_, _ = p.cha8Rand.Read(buf)
	return buf
}
