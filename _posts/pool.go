package main

import (
	"bytes"
	"sync"
)

// BytesBufPool is a type safe wrapper around sync.Pool for
// use with *bytes.Buffer objects.
type BytesBufPool struct {
	pool *sync.Pool
}

// NewBytesBufPool instantiates a new pool of *bytes.Buffer
// objects.
func NewBytesBufPool() *BytesBufPool {
	return &BytesBufPool{
		pool: &sync.Pool{
			New: func() interface{} { return new(bytes.Buffer) },
		},
	}
}

// Get retrieves a *bytes.Buffer from the pool in a type safe manner.
func (b *BytesBufPool) Get() *bytes.Buffer {
	return b.pool.Get().(*bytes.Buffer)
}

// Put returns a *bytes.Buffer object to the pool.
func (b *BytesBufPool) Put(buf *bytes.Buffer) {
	b.pool.Put(buf)
}
