package core

import (
	"runtime"
	"strconv"
	"strings"
)

// @title :  Gid 获取GoroutineId
// @remark:  尽量避免使用基于GoroutineId实现“Goroutine Local Storage”
//	         原因见 https://blog.csdn.net/sb___itfk/article/details/107862104
func Gid() (gid int) {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	gid, err := strconv.Atoi(idField)
	if err != nil {
		gid = -1
	}
	return gid
}
