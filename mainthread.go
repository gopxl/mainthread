package mainthread

import (
	"errors"
	"os"
	"runtime"
)

// CallQueueCap is the capacity of the call queue. This means how many calls to CallNonBlock will not
// block until some call finishes.
//
// The default value is 16 and should be good for 99% usecases.
var CallQueueCap = 16

var (
	callQueue chan func()
	respChan  chan interface{}
)

func init() {
	runtime.LockOSThread()
}

func checkRun() {
	if callQueue == nil {
		panic(errors.New("mainthread: did not call Run"))
	}
}

// Run enables mainthread package functionality. To use mainthread package, put your main function
// code into the run function (the argument to Run) and simply call Run from the real main function.
func Run(run func()) {
	callQueue = make(chan func(), CallQueueCap)
	respChan = make(chan interface{})

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
			os.Exit(0)
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
	callQueue <- func() {
		f()
		respChan <- struct{}{}
	}
	<-respChan
}

// CallErr queues function f on the main thread and returns an error returned by f.
func CallErr(f func() error) error {
	checkRun()
	callQueue <- func() {
		respChan <- f()
	}
	err := <-respChan
	if err != nil {
		return err.(error)
	}
	return nil
}

// CallVal queues function f on the main thread and returns a value returned by f.
func CallVal(f func() interface{}) interface{} {
	checkRun()
	callQueue <- func() {
		respChan <- f()
	}
	return <-respChan
}
