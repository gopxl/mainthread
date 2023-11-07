package mainthread

import (
	"errors"
	"runtime"
)

// CallQueueCap is the capacity of the call queue. This means how many calls to CallNonBlock will not
// block until some call finishes.
//
// The default value is 16 and should be good for 99% usecases.
var CallQueueCap = 16

var (
	callQueue chan func()
	donePool  chan chan any
)

func init() {
	runtime.LockOSThread()
}

func checkRun() {
	if callQueue == nil {
		panic(errors.New("mainthread: did not call Run"))
	}
}

func returnDone(done chan any) {
	donePool <- done
}

// Run enables mainthread package functionality. To use mainthread package, put your main function
// code into the run function (the argument to Run) and simply call Run from the real main function.
//
// Run returns when run (argument) function finishes.
func Run(run func()) {
	callQueue = make(chan func(), CallQueueCap)

	donePool = make(chan chan any, CallQueueCap)
	for i := 0; i < CallQueueCap; i++ {
		donePool <- make(chan any)
	}

	done := make(chan struct{})
	go func() {
		run()
		done <- struct{}{}
	}()

	for {
		select {
		case f := <-callQueue:
			f()
		case <-done:
			return
		}
	}
}

// CallNonBlock queues function f on the main thread and returns immediately. Does not wait until f
// finishes.
func CallNonBlock(f func()) {
	checkRun()
	callQueue <- f
}

// Call queues function f on the main thread and blocks until the function f finishes.
func Call(f func()) {
	checkRun()
	done := <-donePool
	defer returnDone(done)
	callQueue <- func() {
		f()
		done <- struct{}{}
	}
	<-done
}

// CallErr queues function f on the main thread and returns an error returned by f.
func CallErr(f func() error) error {
	return CallVal(f)
}

// CallVal queues function f on the main thread and returns a value returned by f.
func CallVal[T any](f func() T) (t T) {
	checkRun()
	respChan := <-donePool
	defer returnDone(respChan)
	callQueue <- func() {
		respChan <- f()
	}
	v := <-respChan
	if v != nil {
		return v.(T)
	}
	return
}
