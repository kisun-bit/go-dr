//
// syscall相关文档说明：
// 主要向外提供Go调用dll动态链接库，so动态库、lib静态库等等方式

package core

//func CallDLLProc(DLLPath string, proc string, args ...uintptr) (r1, r2 uintptr, err error) {
//	defer func() {
//		if e := recover(); e != nil {
//			r1, r2, err = 0, 0, errors.New(fmt.Sprintf("error in CallDLLProc: %v", e))
//		}
//	}()
//
//	var DLL *syscall.DLL
//	var fProc *syscall.Proc
//
//	DLL, err = syscall.LoadDLL(DLLPath)
//	if err != nil {
//		return 0, 0, err
//	}
//
//	fProc, err = DLL.FindProc(proc)
//	if err != nil {
//		return 0, 0, err
//	}
//
//	r1, r2, _ = fProc.Call(args...)
//	return r1, r2, nil
//}
