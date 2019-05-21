//build +unit
package wgwaiter

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestWgWaiter_2_good_operations(t *testing.T) {

	wg := NewWaiter()

	wg.AddOne()
	go func(wg *WgWaiter) {
		t.Logf("Starting goroutine 1...")
		time.Sleep(1 * time.Second)
		wg.Done()
		t.Logf("Completed goroutine 1")
	}(wg)

	wg.AddOne()
	go func(wg *WgWaiter) {
		t.Logf("Starting goroutine 2...")
		time.Sleep(2 * time.Second)
		wg.Done()
		t.Logf("Completed goroutine 2")
	}(wg)

	err := wg.Wait(5 * time.Second)
	if err != nil {
		t.Errorf("Should be ok, but got %v", err)
	}
}

func TestWgWaiter_2_operations_1_fail(t *testing.T) {
	assert := assert.New(t)
	wg := NewWaiter()

	wg.AddOne()
	go func(wg *WgWaiter) {
		t.Logf("Starting goroutine 1...")
		time.Sleep(1 * time.Second)
		wg.Fail(errors.New("Test Error"))
		t.Logf("Completed goroutine 1")
	}(wg)

	wg.AddOne()
	go func(wg *WgWaiter) {
		t.Logf("Starting goroutine 2...")
		time.Sleep(2 * time.Second)
		wg.Done()
		t.Logf("Completed goroutine 2")
	}(wg)

	err := wg.Wait(5 * time.Second)
	assert.EqualError(err, "Test Error")
}

func TestWgWaiter_timeout(t *testing.T) {
	assert := assert.New(t)
	wg := NewWaiter()

	wg.AddOne()
	go func(wg *WgWaiter) {
		t.Logf("Starting goroutine 1...")
		time.Sleep(1 * time.Second)
		wg.Fail(errors.New("Test Error"))
		t.Logf("Completed goroutine 1")
	}(wg)

	wg.AddOne()
	go func(wg *WgWaiter) {
		t.Logf("Starting goroutine 2...")
		time.Sleep(2 * time.Second)
		wg.Done()
		t.Logf("Completed goroutine 2")
	}(wg)

	err := wg.Wait(200 * time.Millisecond)
	assert.EqualError(err, "Timeout after 200ms")
}

func TestWgWaiter_timeout2(t *testing.T) {
	assert := assert.New(t)
	wg := NewWaiter()

	wg.AddOne()
	go func(wg *WgWaiter) {
		t.Logf("Starting goroutine 1...")
		time.Sleep(3 * time.Second)
		wg.Done()
		t.Logf("Completed goroutine 1")
	}(wg)

	wg.AddOne()
	go func(wg *WgWaiter) {
		t.Logf("Starting goroutine 2...")
		time.Sleep(7 * time.Second)
		wg.Done()
		t.Logf("Completed goroutine 2")
	}(wg)

	err := wg.Wait(5 * time.Second)
	assert.EqualError(err, "Timeout after 5s")
}
