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

func TestWgWaiter_lock(t *testing.T) {
	assert := assert.New(t)
	wg := NewWaiter()
	type data struct {
		one string
		two string
	}
	d := new(data)

	wg.AddOne()
	go func(wg *WgWaiter, dt *data) {
		t.Logf("Starting goroutine 1...")
		wg.Lock()
		dt.one = "one"
		wg.Unlock()
		wg.Done()
		t.Logf("Completed goroutine 1")
	}(wg, d)

	wg.AddOne()
	go func(wg *WgWaiter, dt *data) {
		t.Logf("Starting goroutine 2...")
		time.Sleep(1 * time.Second)
		wg.Lock()
		dt.two = "two"
		wg.Unlock()
		wg.Done()
		t.Logf("Completed goroutine 2")
	}(wg, d)

	err := wg.Wait(2 * time.Second)
	assert.Nil(err)
	assert.Equal(&data{"one", "two"}, d)
}

func TestWgWaiter_lock_no_unlock(t *testing.T) {
	assert := assert.New(t)
	wg := NewWaiter()
	type data struct {
		one string
		two string
	}
	d := new(data)

	wg.AddOne()
	go func(wg *WgWaiter, dt *data) {
		t.Logf("Starting goroutine 1...")
		wg.Lock()
		dt.one = "one"
		wg.Done()
		t.Logf("Completed goroutine 1")
	}(wg, d)

	wg.AddOne()
	go func(wg *WgWaiter, dt *data) {
		t.Logf("Starting goroutine 2...")
		time.Sleep(1 * time.Second)
		wg.Lock()
		dt.two = "two"
		wg.Unlock()
		wg.Done()
		t.Logf("Completed goroutine 2")
	}(wg, d)

	err := wg.Wait(2 * time.Second)
	assert.NotNil(err)
	assert.EqualError(err, "Timeout after 2s")
	assert.Equal(&data{"one", ""}, d)
}
