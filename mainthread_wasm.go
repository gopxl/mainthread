//go:build js && wasm

package mainthread

// Under WASM the entire program runs on the JS event-loop thread. Queuing onto
// a separate main thread would deadlock, so every helper runs the callback
// inline on the calling goroutine.

var CallQueueCap = 16

// Run invokes the supplied function synchronously and returns when it exits.
func Run(run func()) { run() }

// Call runs f inline on the current goroutine.
func Call(f func()) { f() }

// CallNonBlock runs f inline on the current goroutine.
func CallNonBlock(f func()) { f() }

// CallErr runs f inline and returns its error.
func CallErr(f func() error) error { return f() }

// CallVal runs f inline and returns its value.
func CallVal[T any](f func() T) T { return f() }
