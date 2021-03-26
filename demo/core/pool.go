package main

import (
	"fmt"
	"github.com/kisunSea/jpkt/src/core"
	"time"
)

func tHandler(v ...interface{}) error {
	time.Sleep(1 * time.Second)
	//panic("1111")
	fmt.Println(core.Gid(), v)
	return nil
}

func TestPool() {
	pool, _ := core.NewJPool(10)

	for i := 0; i < 20; i++ {
		t := core.NewJTask(tHandler, []interface{}{i, "test"})
		_ = pool.Put(t)
	}

	pool.Start()
	pool.Close()
}


