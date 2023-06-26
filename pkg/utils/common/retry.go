package common

import (
	"math"
	"sync/atomic"
	"time"
	"unsafe"
)

type state [2]time.Duration

type fibonacciBackoff struct {
	state unsafe.Pointer
}
type Backoff interface {
	Next() (next time.Duration, stop bool)
}

func NewFibonacci(base time.Duration) Backoff {
	if base <= 0 {
		panic("base must be greater than 0")
	}

	return &fibonacciBackoff{
		state: unsafe.Pointer(&state{0, base}),
	}
}

func (b *fibonacciBackoff) Next() (time.Duration, bool) {
	for {
		curr := atomic.LoadPointer(&b.state)
		currState := (*state)(curr)
		next := currState[0] + currState[1]

		if next <= 0 {
			return math.MaxInt64, false
		}

		if atomic.CompareAndSwapPointer(&b.state, curr, unsafe.Pointer(&state{currState[1], next})) {
			return next, false
		}
	}
}
