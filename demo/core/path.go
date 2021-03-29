package main

import (
	"github.com/kisunSea/jpkt/src/core"
)

func testPath() {
	core.GetCurrentPath()
	core.IsFile("/home/zk")
	core.IsDir("/home/zk")
	_, _ = core.PathExists("/home/zk")
	core.EnumAllFilePathInDir("/home/zk")
}
