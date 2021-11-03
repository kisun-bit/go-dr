package main

import (
	"bytes"
	"fmt"
	"github.com/kisunSea/go_dr/src/cmd"
	"github.com/kr/pretty"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"io/ioutil"
	"os/exec"
	"time"
)

func ExecuteSyncCmd(cmdLine string) (r int, out string, err error) {
	var output bytes.Buffer

	cmd_ := exec.Command("sh")
	in := bytes.NewBuffer(nil)
	cmd_.Stdin = in
	cmd_.Stdout = &output
	go func() {
		in.WriteString(fmt.Sprintf("%s\n", cmdLine))
	}()

	err = cmd_.Start()
	if err != nil {
		return 1, "", err
	}
	err = cmd_.Wait() // wait
	if err != nil {
		return 1, "", err
	}

	// TODO 编码转换，暂时仅考虑GBK的转换...
	utf8Raw, codexErr := GbkToUtf8(output.Bytes())
	if codexErr != nil {
		return 0, output.String(), nil
	} else {
		return 0, string(utf8Raw), nil
	}
}

func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func ControlRunCmd() {
	// Start a long-running process, capture stdout and stderr
	findCmd := cmd.NewCmd("cmd", "dir")
	statusChan := findCmd.Start() // non-blocking

	ticker := time.NewTicker(2 * time.Second)

	// Print last line of stdout every 2s
	go func() {
		for range ticker.C {
			status := findCmd.Status()
			n := len(status.Stdout)
			fmt.Println(status.Stdout[n-1])
		}
	}()

	// Stop command after 1 hour
	go func() {
		<-time.After(1 * time.Hour)
		_ = findCmd.Stop()
	}()

	// Check if command is done
	select {
	case finalStatus := <-statusChan:
		// done
		fmt.Println("===========:\n", finalStatus)
	default:
		// no, still running
	}

	// Block waiting for command to exit, be stopped, or be killed
	finalStatus := <-statusChan

	fmt.Println(finalStatus)
}

func BlockRunCmd() {
	pretty.Log(cmd.SyncExecShell("ll"))
}

func BlockRunCmdV2() {
	c := cmd.NewCmd("iscsicli --sessionlist")
	s := <-c.Start()
	pretty.Log(s)
}

func main() {
	//BlockRunCmd()
	BlockRunCmdV2()
	//ControlRunCmd()
	//fmt.Println(ExecuteSyncCmd("dir"))  // 0  "docs  example  go.mod  go.sum  src  test" <nil>
}
