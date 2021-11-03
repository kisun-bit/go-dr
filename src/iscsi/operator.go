package iscsi

import "runtime"

type Common interface {
	// Check that the IQN target is connected
	CheckTarget() (err error)
	// Register server IQN
	AddTargets() (err error)
	// After logging in to IQN, you can map local virtual disks
	LoginTarget() (err error)
	// Logs out the device from the specified IQN
	LogoutTarget() (err error)
	// Connects to the specified IQN
	ConnectTarget() (err error)
	// TODO more...
}

type TargetInfo struct {
	IP   string `json:"ip"`
	Port string `json:"port"`
	IQN  string `json:"iqn"`

	// The following attributes are specific to Windows

	Initiator string
}

func NewOperator(ip, port, iqn string) Common {

	ti := TargetInfo{IP: ip, Port: port, IQN: iqn}

	switch runtime.GOOS {
	case "linux":
		return NewISCSIOpLinux(ti)
	case "windows":
		return NewISCSIHelper()
	default:
		return nil
	}
}
