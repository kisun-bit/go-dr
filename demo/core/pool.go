package main

import (
	"fmt"
	"jpkt/src/core"
	"sync"
	"sync/atomic"
	"time"
)

func TestNewPool() {
	pool := core.NewPool(1000, 10000)
	defer pool.Release()

	iterations := 1000000
	var counter uint64

	wg := sync.WaitGroup{}
	wg.Add(iterations)
	for i := 0; i < iterations; i++ {
		arg := uint64(1)
		job := func() {
			defer wg.Done()
			atomic.AddUint64(&counter, arg)
		}

		pool.JobQueue <- job
	}
	wg.Wait()

	counterFinal := atomic.LoadUint64(&counter)
	if uint64(iterations) != counterFinal {
		fmt.Println(fmt.Errorf("iterations %v is not equal counterFinal %v", iterations, counterFinal))
	}
}

const (
	runTimes  = 10000
	poolSize  = 500
	queueSize = 50
	N         = 1
)

func demoTask(i, j int) {
	time.Sleep(time.Millisecond * 10)
	fmt.Println(fmt.Sprintf("%d-%d", i, j))
}

func TGoroutine() {
	var wg sync.WaitGroup
	for i := 0; i < N; i++ {
		wg.Add(runTimes)

		for j := 0; j < runTimes; j++ {
			go func() {
				defer wg.Done()
				demoTask(i, j)
			}()
		}

		wg.Wait()
	}
}

func TGpool() {
	pool := core.NewPool(poolSize, queueSize)
	defer pool.Release()
	var wg sync.WaitGroup

	for i := 0; i < N; i++ {
		wg.Add(runTimes)
		for j := 0; j < runTimes; j++ {
			pool.JobQueue <- func() {
				defer wg.Done()
				demoTask(i, j)
			}
		}
		wg.Wait()
	}
}

func main() {
	TestNewPool()
	_ = TGoroutine
	TGpool()
}
