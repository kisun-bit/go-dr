package main

import (
	"fmt"
	"github.com/kisun-bit/go_dr/src/core"
)

func testXXHash() {
	//hashcode, err := core.XXHashFileWithLowCollision(`D:\download\ubuntu-18.04.5-desktop-amd64.iso`, 1024)
	hashcode, err := core.XXHashFileWithLowCollision(`D:\tmp\WeChat\WeChat Files\wxid_4zz11sd1b9y022\FileStorage\File\2021-04\rongan-cd.exe`, 1024)
	if err != nil {
		print(err)
		return
	}
	fmt.Println("hashcode:\t", hashcode)
}

//func main() {
//	start := time.Now().Nanosecond()
//	testXXHash()
//	fmt.Printf("总耗时: %d", time.Now().Nanosecond()-start)
//}
