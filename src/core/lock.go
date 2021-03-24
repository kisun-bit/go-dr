package core

import (
	"sync/atomic"
)

// 自旋锁
//
// 注意：
//    1. 请在明确知道“代码段执行时长很短的情况下”使用；
//	  2. 单核CPU上要尽量避免使用自旋锁；
//    3. 该版本的自旋锁进行了优化，避免了空耗CPU的情况，底层调用了“PAUSE”指令；
type SpinLock int32

func (sl *SpinLock) Lock() {
	for !atomic.CompareAndSwapInt32((*int32)(sl), 0, 1) {
	}
}

func (sl *SpinLock) UnLock() {
	atomic.StoreInt32((*int32)(sl), 0)
}

// 文件锁（写锁）
// 注意：
//    1. 该锁支持跨进程实现互斥
//    2. 场景见于串行化执行一组进程，例如：命令行执行配置防火墙的一组命令
//    3. win平台下不支持, 仅支持linux
//type FileLock struct {
//	fPath string
//	f     *os.File
//}
//
//func New(fPath string) *FileLock {
//	return &FileLock{
//		fPath: fPath,
//	}
//}
//
//func (fL *FileLock) Lock() error {
//	f, err := os.Open(fL.fPath)
//	if err != nil {
//		return err
//	}
//	fL.f = f
//	err = syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB)
//	if err != nil {
//		return fmt.Errorf("cannot flock directory %s - %s", fL.fPath, err)
//	}
//	return nil
//}
//
//func (fL *FileLock) Unlock() error {
//	defer fL.f.Close()
//	return syscall.Flock(int(fL.f.Fd()), syscall.LOCK_UN)
//}
