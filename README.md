# wgwaiter
If you need a WaitGroup with timeout and error management, then this package may be what you are looking for.

It eliminates some routine with channels, mutex and timeout when running multiple goroutinces. 
It is an alternative to [Go's context](https://blog.golang.org/context) and [errgroup](https://godoc.org/golang.org/x/sync/errgroup)

# Usage
#### Simple case
```go
package main

import (
    "github.com/kipriz/wgwaiter"
    "log"
    "time"
)

func main() {
    wg := wgwaiter.NewWaiter()

    wg.AddOne()
    go func(wg *wgwaiter.WgWaiter) {
        log.Println("Starting goroutine 1...")
        time.Sleep(1 * time.Second)
        wg.Done() //important to call Done or Fail methods
        log.Println("Completed goroutine 1")
    }(wg)

    wg.AddOne()
    go func(wg *wgwaiter.WgWaiter) {
        log.Println("Starting goroutine 2...")
        time.Sleep(2 * time.Second)
        wg.Done() //important to call Done or Fail methods
        log.Println("Completed goroutine 2")
    }(wg)

    err := wg.Wait(5 * time.Second)
    if err != nil {
        log.Fatalln("Should be ok, but got %v", err)
    }
}
``` 

#### Locks
```go
func main() {
    wg := NewWaiter()
    type data struct {
        one string
        two string
    }
    d := new(data)
    
    wg.AddOne()
    go func(wg *WgWaiter, dt *data) {
        wg.Lock()
        dt.one = "one"
        wg.Unlock()
        wg.Done()
    }(wg, d)
    
    wg.AddOne()
    go func(wg *WgWaiter, dt *data) {
        time.Sleep(1 * time.Second)
        wg.Lock()
        dt.two = "two"
        wg.Unlock()
        wg.Done()
    }(wg, d)
    
    err := wg.Wait(2 * time.Second)
}
```