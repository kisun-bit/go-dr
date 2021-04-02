/* @Title
     标准化 “命令执行”
   @Description
     命令执行的结果有三种：执行结果码、执行输出、错误输出
     命令执行的形式可分为下述最主要的三种
     1. 一次执行完毕，并等待linux系统返回的执行结果；
     2. 执行中需要输入参数才能完成命令，可能涉及多次输入；
     3. 仅调用执行即可，不关心命令的返回；
     ...
   @Remark
*/
package command

import "runtime"

type Cmd interface {
	// 执行命令行并返回结果
	// args: 命令行参数
	// return: 进程的pid, 命令行结果, 错误消息
	Exec(args ...string) (int, string, error)

	// 异步执行命令行并通过channel返回结果
	// stdout: chan结果
	// args: 命令行参数
	// return: 进程的pid
	// exception: 协程内的命令行发生错误时,会panic异常
	ExecAsync(stdout chan string, args ...string) int

	// 执行命令行(忽略返回值)
	// args: 命令行参数
	// return: 错误消息
	ExecIgnoreResult(args ...string) error
}

// Command的初始化函数
func NewCommand() Cmd {
	var cmd Cmd

	switch runtime.GOOS {
	case "linux":
		cmd = NewLinuxCmd()
	case "windows":
		cmd = NewWindowsCmd()
	default:
		cmd = NewLinuxCmd()
	}

	return cmd
}
