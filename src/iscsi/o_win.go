package iscsi

import (
	"errors"
	"fmt"
	"github.com/kisun-bit/go_dr/src/cmd"
	"strings"
)

type OpWin struct {
	TI TargetInfo
}

func NewISCSIHelper() *OpWin {
	return new(OpWin)
}

func (i *OpWin) CheckTarget(target string) (err error) {
	var out string

	if _, out, err = cmd.SyncExecBin("iscsicli listTargets"); err != nil {
		logger.Fmt.Errorf("listTargets failed: %s", err)
	}

	logger.Fmt.Errorf("CheckTarget output targets is:\n %s", out)

	if strings.Contains(out, target) {
		logger.Fmt.Debugf("CheckTarget is ok. %s", target)
		return nil
	}

	logger.Fmt.Debugf("check target is failed: %s", target)
	return errors.New("target is not found")
}

func (i *OpWin) AddTargets(ip, port string) (err error) {

	if _, _, err = cmd.SyncExecBin(fmt.Sprintf("iscsicli addTargetPortal %s %s", ip, port)); err != nil {
		logger.Fmt.Errorf("AddTargets failed: %s", err)
		return err
	}

	logger.Fmt.Debugf("AddTargets is ok")
	return nil
}

func (i *OpWin) LoginTarget(target string) (err error) {
	if _, _, err = cmd.SyncExecBin(
		fmt.Sprintf(
			"iscsicli persistentlogintarget %s T * * * * * * * * * * * * * * * 0",
			strings.ToLower(target))); err != nil {
		logger.Fmt.Errorf("LoginTarget failed: %s\n", err)
	}

	if _, _, err = cmd.SyncExecBin(
		fmt.Sprintf(
			"iscsicli logintarget %s T * * * * * * * * * * * * * * * 0",
			strings.ToLower(target))); err != nil {
		logger.Fmt.Errorf("persistentlogintarget failed: %s\n", err)
		return err
	}

	logger.Fmt.Debugf("LoginTarget is ok. %s\n", target)
	return nil
}

func (i *OpWin) DeletePersistentDrive(vol string) (err error) {
	if _, _, err = cmd.SyncExecBin(
		fmt.Sprintf(
			"iscsicli RemovePersistentDevice %s:\\", strings.ToUpper(vol))); err != nil {
		logger.Fmt.Debugf("DeletePersistentDrive failed: %s\n", err)
		return err
	}
	return nil
}

func (i *OpWin) AddPersistentDrive(vol string) (err error) {
	var out string
	if _, out, err = cmd.SyncExecBin(
		fmt.Sprintf(
			`iscsicli AddPersistentDevice %s:\`, vol)); !strings.Contains(out, "成功") {
		logger.Fmt.Errorf("BindPersistentDevice execute cmd failed: %s", out)
		return fmt.Errorf("bind persistent device %s:\\ failed", vol)
	}
	return nil
}

func (i *OpWin) RemovePersistentTarget(initiator, target, port, tIP, tPort string) (err error) {
	if _, _, err = cmd.SyncExecBin(
		fmt.Sprintf("iscsicli RemovePersistentTarget %s %s %s %s %s",
			initiator, target, port, tIP, tPort)); err != nil {
		logger.Fmt.Debugf("RemovePersistentTarget failed: %s\n", err)
		return err
	}
	return nil
}

func (i *OpWin) LogoutTarget(sessionID string) (err error) {
	if _, _, err = cmd.SyncExecBin(
		fmt.Sprintf(
			"iscsicli logoutTarget %s",
			strings.ToLower(sessionID))); err != nil {
		logger.Fmt.Debugf("logoutTarget failed: %s\n", err)
		return err
	}

	logger.Fmt.Debugf("LogoutTarget is ok. %s\n", sessionID)
	return nil
}

func (i *OpWin) ConnectTarget(ip, port, targetIQN string) (err error) {
	logger.Fmt.Debugf("connect targetIQN: %s:%s -> %s\n", ip, port, targetIQN)

	if err = i.AddTargets(ip, port); err != nil {
		return err
	}

	//if err = i.CheckTarget(targetIQN); err != nil {
	//	//return nil  // 暂时忽略
	//}

	if err = i.LoginTarget(targetIQN); err != nil {
		//return err // TODO 暂时忽略错误
	}

	logger.Fmt.Debugf("connenct targetIQN is ok: %s\n", targetIQN)
	return nil
}
