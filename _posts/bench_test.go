package main

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
)

func fmtStatStr(stat string, tags map[string]string) string {
	parts := []string{}
	for k, v := range tags {
		if v != "" {
			parts = append(parts, fmt.Sprintf("%s:%s", k, v))
		}
	}

	return fmt.Sprintf("%s|%s", stat, strings.Join(parts, ","))
}

func fmtStatStrAfter(stat string, tags map[string]string) string {
	b := bytes.NewBufferString(stat)
	b.WriteString("|")

	first := true
	for k, v := range tags {
		if v != "" {
			if !first {
				// do not append a ',' at the
				// beginning of the first iteration.
				b.WriteString(",")
			}

			b.WriteString(k)
			b.WriteString(":")
			b.WriteString(v)

			first = false // automatically set to falsey after first iteration
		}
	}

	return b.String()
}

var test = struct {
	stat string
	tags map[string]string
}{
	"handler.sample",
	map[string]string{
		"os":     "ios",
		"locale": "en-US",
	},
}

func Benchmark_fmtStatStrBefore(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmtStatStr(test.stat, test.tags)
	}
}

func Benchmark_fmtStatStrAfter(b *testing.B) {
	for i := 0; i < b.N; i++ {
		fmtStatStrAfter(test.stat, test.tags)
	}
}
