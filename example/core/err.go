package main

import (
	"errors"
	"fmt"
	"github.com/kisunSea/go_dr/src/core"
	"github.com/kisunSea/go_dr/src/datahandle"
	"github.com/kisunSea/go_dr/src/meta"
)

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
		meta.ErrUnknown,
		"UnknownError",
		fmt.Sprintf("内部异常，代码[%v]", datahandle.FmtErrCode2String(meta.ErrUnknown)),
		"Can't do this thing...")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println(err.ErrorDetail())
		err.AddMoreDebug("这是新增调试信息1")
		err.AddMoreDebug("这是新增调试信息2")
		err.AddMoreDebug("这是新增调试信息3")
		err.ChangeDescription("修改描述信息1")
		fmt.Println(err.ErrorDetail())
	}

	err2 := errors.New("EEEEEEEEEEEEEEE")
	se := core.StandardizeErr(err2)
	fmt.Println(se.ErrorDetail())
}

func test5() {
	defer core.CatchPanicErr(nil)
	panic(core.RaiseStandardError(
		meta.ErrUnknown,
		"UnknownError",
		fmt.Sprintf("内部异常，代码[%v]", datahandle.FmtErrCode2String(meta.ErrUnknown)),
		"Can't do this thing..."))
}

func test6() {
	core.StandardPanic(meta.ErrNoFreePort, "Test", "test...test...test", "debug", nil)
}

//func main() {
//	defer core.CatchPanicErr(nil)
//	test1()
//	test5()
//	test6()
//}
