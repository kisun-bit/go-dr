package main

import (
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/kisunSea/jpkt/src/core"
	"github.com/kisunSea/jpkt/src/datahandle"
	"github.com/kisunSea/jpkt/src/meta"
)

var Lg *zap.SugaredLogger

var Lg2 *zap.SugaredLogger

//func init() {
//	Lg = log.GetJLoggerByConf(`D:\workspace\Jpkt\demo\log`, "log", "default")
//	defer Lg.Sync()
//
//	Lg2 = log.GetJLoggerByMapConf(
//		map[string]interface{}{
//			"Logger":       "default",
//			"Filename":     "./demo/log/jpkt.log",
//			"MaxSize":      30,
//			"MaxBackups":   3,
//			"MaxAge":       30,
//			"Compress":     false,
//			"Level":        "debug",
//			"SplitByLevel": false,
//		},
//	)
//	defer Lg2.Sync()
//}

func test1() {
	test2()
}

func test2() {
	test3()
}

func test3() {
	test4()
}

func test4() {
	err := core.RaiseStandardError(
		meta.JErrUnknown,
		"UnknownError",
		fmt.Sprintf("内部异常，代码[%v]", datahandle.FmtErrCode2String(meta.JErrUnknown)),
		"Can't do this thing...")
	if err != nil {
		Lg.Info(err.Error())
		Lg.Info(err.ErrorDetail())
		err.AddMoreDebug("这是新增调试信息1")
		err.AddMoreDebug("这是新增调试信息2")
		err.AddMoreDebug("这是新增调试信息3")
		err.ChangeDescription("修改描述信息1")
		Lg.Info(err.ErrorDetail())
	}

	err2 := errors.New("EEEEEEEEEEEEEEE")
	se := core.StandardizeErr(err2)
	Lg.Info(se.ErrorDetail())
}

func test5() {
	defer core.CatchPanicErr(Lg)
	panic(core.RaiseStandardError(
		meta.JErrUnknown,
		"UnknownError",
		fmt.Sprintf("内部异常，代码[%v]", datahandle.FmtErrCode2String(meta.JErrUnknown)),
		"Can't do this thing..."))
}

func test6() {
	core.StandardPanic(meta.JErrNoFreePort, "Test", "test...test...test", "debug", nil)
}

func testLog() {
	Aaa()
}

func testLogDemo() {
	testLog()
	go test1()
	test4()
	test5()
	test6()
	time.Sleep(5 * 1000)
}

//func main() {
//	Lg2.Info("hhhhhhhhhhhhhhhhh")
//}
