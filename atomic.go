// Copyright (c) 2022. Homeland Interactive Technology Ltd. All rights reserved.

package ruixuego

import (
	"sync/atomic"
)

// Bool is an atomic Boolean.
type Bool struct{ v uint32 }

// NewBool creates a Bool.
func NewBool(initial bool) *Bool {
	return &Bool{boolToInt(initial)}
}

// Load atomically loads the Boolean.
func (b *Bool) Load() bool {
	return truthy(atomic.LoadUint32(&b.v))
}

// CAS is an atomic compare-and-swap.
func (b *Bool) CAS(old, new bool) bool {
	return atomic.CompareAndSwapUint32(&b.v, boolToInt(old), boolToInt(new))
}

// Store atomically stores the passed value.
func (b *Bool) Store(new bool) {
	atomic.StoreUint32(&b.v, boolToInt(new))
}

// Swap sets the given value and returns the previous value.
func (b *Bool) Swap(new bool) bool {
	return truthy(atomic.SwapUint32(&b.v, boolToInt(new)))
}

// Toggle atomically negates the Boolean and returns the previous value.
func (b *Bool) Toggle() bool {
	return truthy(atomic.AddUint32(&b.v, 1) - 1)
}

func truthy(n uint32) bool {
	return n&1 == 1
}

func boolToInt(b bool) uint32 {
	if b {
		return 1
	}
	return 0
}
