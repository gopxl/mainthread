//go:build js && wasm

package mainthread

// Under WASM the entire program runs on the JS event-loop thread. Queuing onto
// a separate main thread would deadlock, so the blocking helpers (Run, Call,
// CallErr, CallVal) run the callback inline. CallNonBlock uses a goroutine to
// preserve its "return immediately" semantic.

// CallQueueCap exists for source compatibility with the desktop build; the
// WASM shim does not queue anything, so the value is unused.
var CallQueueCap = 16

// Run invokes the supplied function synchronously and returns when it exits.
// On WASM there is no separate main thread to pump, so Run is effectively a
// direct call.
func Run(run func()) { run() }

// Call runs f inline on the current goroutine. Equivalent to the desktop
// behavior of queueing onto the main thread and blocking until done, because
// under WASM the calling goroutine is already on the main (and only) thread.
func Call(f func()) { f() }

// CallNonBlock runs f on a new goroutine so the caller is not blocked. This
// matches the desktop semantic ("queue and return") as closely as possible
// under a single-threaded runtime; f will execute cooperatively once the
// caller yields to the scheduler.
func CallNonBlock(f func()) { go f() }

// CallErr runs f inline and returns its error.
func CallErr(f func() error) error { return f() }

// CallVal runs f inline and returns its value.
func CallVal[T any](f func() T) T { return f() }
