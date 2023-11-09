package mainthread_test

import (
	"errors"
	"io"
	"testing"

	"github.com/gopxl/mainthread/v2"
	"github.com/stretchr/testify/require"
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

func TestCallVal_Resp(t *testing.T) {
	run := func() {
		f := func() *int {
			v := 1
			return &v
		}
		v := mainthread.CallVal(f)
		require.NotNil(t, v)
		require.Equal(t, 1, *v)
	}
	mainthread.Run(run)
}

func TestCallVal_nilResp(t *testing.T) {
	run := func() {
		f := func() *int {
			return nil
		}
		v := mainthread.CallVal(f)
		require.Nil(t, v)
	}
	mainthread.Run(run)
}

func TestCallErr_Resp(t *testing.T) {
	run := func() {
		f := func() error {
			return io.EOF
		}
		v := mainthread.CallErr(f)
		require.NotNil(t, v)
		require.EqualError(t, v, io.EOF.Error())
	}
	mainthread.Run(run)
}

func TestCallErr_nilResp(t *testing.T) {
	run := func() {
		f := func() error {
			return nil
		}
		v := mainthread.CallErr(f)
		require.NoError(t, v)
	}
	mainthread.Run(run)
}
