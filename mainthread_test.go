package mainthread_test

import (
	"errors"
	"testing"

	"github.com/gopxl/mainthread/v2"
)

func BenchmarkCall(b *testing.B) {
	run := func() {
		f := func() {}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mainthread.Call(f)
		}
	}
	mainthread.Run(run)
}

func BenchmarkCallErr(b *testing.B) {
	run := func() {
		f := func() error {
			return errors.New("foo")
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mainthread.CallErr(f)
		}
	}
	mainthread.Run(run)
}

func BenchmarkCallVal(b *testing.B) {
	run := func() {
		f := func() int {
			return 42
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			mainthread.CallVal(f)
		}
	}
	mainthread.Run(run)
}
