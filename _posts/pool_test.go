package main

import (
	"bytes"
	"testing"
)

func Benchmark_noPool(b *testing.B) {
	for i := 0; i < b.N; i++ {
		var buf bytes.Buffer
		buf.Write([]byte("Hello"))
	}
}

func Benchmark_pool(b *testing.B) {
	p := NewBytesBufPool()
	for i := 0; i < b.N; i++ {
		buf := p.Get()
		buf.Write([]byte("Hello"))
		p.Put(buf)
	}
}
