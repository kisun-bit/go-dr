package command

import (
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"
)

type WindowsCmd struct {
}

// WindowsCommand的初始化函数
func NewWindowsCmd() *WindowsCmd {
	return &WindowsCmd{}
}

// 执行命令行并返回结果
// args: 命令行参数
// return: 进程的pid, 命令行结果, 错误消息
func (wc *WindowsCmd) Exec(args ...string) (int, string, error) {
	args = append([]string{"-c"}, args...)
	cmd := exec.Command("cmd", args...)

	cmd.SysProcAttr = &syscall.SysProcAttr{}

	outpip, err := cmd.StdoutPipe()
	if err != nil {
		return 0, "", err
	}
	defer outpip.Close()

	err = cmd.Start()
	if err != nil {
		return 0, "", err
	}

	out, err := ioutil.ReadAll(outpip)
	if err != nil {
		return 0, "", err
	}

	return cmd.Process.Pid, string(out), nil
}

// 异步执行命令行并通过channel返回结果
// stdout: chan结果
// args: 命令行参数
// return: 进程的pid
// exception: 协程内的命令行发生错误时,会panic异常
func (wc *WindowsCmd) ExecAsync(stdout chan string, args ...string) int {
	var pidChan = make(chan int, 1)

	go func() {
		args = append([]string{"-c"}, args...)
		cmd := exec.Command("cmd", args...)

		cmd.SysProcAttr = &syscall.SysProcAttr{}

		outpip, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}
		defer outpip.Close()

		err = cmd.Start()
		if err != nil {
			panic(err)
		}

		pidChan <- cmd.Process.Pid

		out, err := ioutil.ReadAll(outpip)
		if err != nil {
			panic(err)
		}

		stdout <- string(out)
	}()

	return <-pidChan
}

// 执行命令行(忽略返回值)
// args: 命令行参数
// return: 错误消息
func (wc *WindowsCmd) ExecIgnoreResult(args ...string) error {
	args = append([]string{"-c"}, args...)
	cmd := exec.Command("cmd", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{}

	err := cmd.Run()

	return err
}
