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
package core

func DelayExecuteSh() {
	// TODO 延迟执行命令
}

func ExecuteSh() {
	// TODO 执行命令，添加超时、输入、编码、等
}

type ASyncProcessSh struct {
	// TODO 异步进程类型，不必等待命令执行的返回，
}
