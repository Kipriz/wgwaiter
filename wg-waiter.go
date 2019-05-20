package wgwaiter

import (
	"sync"
	"time"
)

type WgWaiter struct {
	wg         sync.WaitGroup
	errChannel chan error
	finished   chan bool
}

func NewWaiter() *WgWaiter {
	return &WgWaiter{
		wg:         sync.WaitGroup{},
		errChannel: make(chan error),
		finished:   make(chan bool, 1),
	}
}

func (wg *WgWaiter) AddOne() {
	wg.wg.Add(1)
}

func (wg *WgWaiter) Add(delta int) {
	wg.wg.Add(delta)
}

func (wg *WgWaiter) Done() {
	wg.wg.Done()
}

func (wg *WgWaiter) Fail(err error) {
	wg.errChannel <- err
	wg.wg.Done()
}

func (wg *WgWaiter) Wait(timeout time.Duration) error {
	go func() {
		wg.wg.Wait()
		wg.finished <- true
	}()
	var err error
	select {
	case err = <-wg.errChannel:
		return err
	case <-wg.finished:
		return nil
	case <-time.After(timeout):
		return TimeoutError(timeout.String())
	}
}

type TimeoutError string

func (te TimeoutError) Error() string {
	return "Timeout after " + string(te)
}
