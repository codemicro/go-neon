package tests

import (
	"fmt"
	"github.com/codemicro/go-neon/examples/benchmark/templates"
	"testing"
)

func BenchmarkNeonTemplate1(b *testing.B) {
	benchmarkNeonTemplate(b, 1)
}

func BenchmarkNeonTemplate10(b *testing.B) {
	benchmarkNeonTemplate(b, 10)
}

func BenchmarkNeonTemplate100(b *testing.B) {
	benchmarkNeonTemplate(b, 100)
}

func benchmarkNeonTemplate(b *testing.B, rowsCount int) {
	rows := getBenchRows(rowsCount)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			templates.BenchPage(rows)
		}
	})
}

func getBenchRows(n int) []templates.BenchRow {
	rows := make([]templates.BenchRow, n)
	for i := 0; i < n; i++ {
		rows[i] = templates.BenchRow{
			ID:      i,
			Message: fmt.Sprintf("message %d", i),
			Print:   ((i & 1) == 0),
		}
	}
	return rows
}
