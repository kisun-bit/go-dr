package main

import (
	"github.com/kisun-bit/go_dr/src/datahandle"
)

func testPath() {
	datahandle.GetCurrentPath()
	datahandle.IsFile("/home/zk")
	datahandle.IsDir("/home/zk")
	_, _ = datahandle.PathExists("/home/zk")
	datahandle.EnumAllFilePathInDir("/home/zk")
}
