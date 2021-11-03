package iscsi

import (
	"errors"
	"fmt"
	"github.com/kisunSea/go_dr/src/cmd"
	"strings"
	"time"
)

type OpLinux struct {
	TI TargetInfo
}

func NewISCSIOpLinux(t TargetInfo) *OpLinux {
	return &OpLinux{TI: t}
}

func (i *OpLinux) CheckTarget() (err error) {
	var out string

	if _, out, err = cmd.SyncExecBin(fmt.Sprintf("lsscsi -t | grep \"%s,\" | awk '{print $NF}'", i.TI.IQN),
	); err != nil {
		logger.Fmt.Errorf("CheckTarget failed: %s", err)
	}

	logger.Fmt.Debugf("CheckTarget output targets is:\n %s", out)

	if strings.TrimSpace(out) != "" {
		logger.Fmt.Debugf("CheckTarget is ok. %s", i.TI.IQN)
		return nil
	}

	logger.Fmt.Errorf("check target is failed: %s", i.TI.IQN)
	return errors.New("target is not found")
}

func (i *OpLinux) AddTargets() (err error) {

	if _, _, err = cmd.SyncExecBin(fmt.Sprintf("iscsiadm -m discovery -t st -p %s:%s", i.TI.IP, i.TI.Port),
	); err != nil {
		logger.Fmt.Errorf("AddTargets failed: %s", err)
		return err
	}

	logger.Fmt.Debugf("AddTargets is ok")
	return nil
}

func (i *OpLinux) LoginTarget() (err error) {
	var out string
	if _, _, err = cmd.SyncExecBin(fmt.Sprintf("iscsiadm -m node -T %s -l", strings.ToLower(i.TI.IQN)),
	); err != nil {
		logger.Fmt.Errorf("LoginTarget failed: %s | %s\n", out, err)
		return err
	}

	logger.Fmt.Debugf("LoginTarget is ok. %s\n", i.TI.IQN)
	return nil
}

func (i *OpLinux) LogoutTarget() (err error) {
	var out string

	if _, _, err = cmd.SyncExecBin(
		fmt.Sprintf(
			"iscsiadm -m node -T %s -u",
			strings.ToLower(i.TI.IQN))); err != nil {
		logger.Fmt.Warnf("logoutTarget failed: %s | %s", out, err)
		return err
	}

	_, _, _ = cmd.SyncExecBin(fmt.Sprintf("iscsiadm -m node -o delete –T  %s -p %s:%s",
		i.TI.IQN, i.TI.IP, i.TI.Port))

	logger.Fmt.Debugf("LogoutTarget is ok. `%s`", i.TI.IQN)
	return nil
}

func (i *OpLinux) RefreshSession() (err error) {
	if _, _, err = cmd.SyncExecBin("iscsiadm -m session –R"); err != nil {
		logger.Fmt.Errorf("RefreshSession failed: %s\n", err)
		return err
	}
	logger.Fmt.Debugf("RefreshSession is ok.")
	return nil
}

func (i *OpLinux) ConnectTarget() (err error) {
	logger.Fmt.Debugf("connect targetTI.IQN: %s:%s -> %s", i.TI.IP, i.TI.Port, i.TI.IQN)

	if err = i.AddTargets(); err != nil {
		// return err
	}

	if err = i.LoginTarget(); err != nil {
		return err
	}

	time.Sleep(5 * time.Second)

	if err = i.CheckTarget(); err != nil {
		return err
	}

	logger.Fmt.Debugf("connenct targetTI.IQN is ok: %s", i.TI.IQN)
	return nil
}
