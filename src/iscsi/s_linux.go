package iscsi

import (
	"errors"
	"fmt"
	"github.com/kisunSea/go_dr/src/cmd"
	"regexp"
	"strings"
)

/////////////////////////////////////////////////////////////////

var (
	iSCSILine1Pattern = regexp.MustCompile("Host:(?P<host>.*?)Channel:(?P<chanel>.*?)Target:(?P<target>.*?)Lun:(?P<lun>.*?)$")
	iSCSILine2Pattern = regexp.MustCompile("Vendor:(?P<vendor>.*?)Model:(?P<model>.*?)Rev:(?P<rev>.*?)$")
	iSCSILine3Pattern = regexp.MustCompile("Type:(?P<vendor>.*?)ANSI SCSI revision:(?P<revision>.*?)$")
)

type ISCSIDriver struct {
	Host      string
	Channel   string
	Target    string
	Lun       string
	UUID      string
	Vendor    string
	Model     string
	Rev       string
	Type      string
	Revision  string
	DrivePath string
	IQNInfo   string
}

func __fixedNo(str string) string {
	str = strings.TrimSpace(str)
	if str == "00" {
		return "0"
	} else if strings.HasPrefix(str, "0") {
		return str[1:]
	} else {
		return str
	}
}

func GetISCSIDrivers() (ids []ISCSIDriver, err error) {
	var out string
	var r int

	for i := 0; i < 10; i++ {
		r, out, err = cmd.SyncExecBin("lsscsi -c")
		if out != "" {
			break
		}
	}

	if out == "" {
		return nil, errors.New("exec lsscsi failed")
	}

	lines := strings.Split(out, "\n")
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		if line = strings.TrimSpace(line); line == "" {
			continue
		}
		if strings.HasPrefix(line, "Host") && (i+2) <= len(lines) {
			match1 := iSCSILine1Pattern.FindStringSubmatch(line)
			match2 := iSCSILine2Pattern.FindStringSubmatch(lines[i+1])
			match3 := iSCSILine3Pattern.FindStringSubmatch(lines[i+2])
			if len(match1) < 5 || len(match2) < 4 || len(match3) < 3 {
				continue
			}

			uuid := strings.ReplaceAll(
				fmt.Sprintf("[%v:%v:%v:%v]",
					__fixedNo(match1[1]),
					__fixedNo(match1[2]),
					__fixedNo(match1[3]),
					__fixedNo(match1[4])),
				"scsi", "")

			iqnInfo := ""
			if r, iqnInfo, err = cmd.SyncExecBin(
				fmt.Sprintf(`lsscsi -t | grep "\[%v\]"`, uuid[1:len(uuid)-1])); err != nil || r != 0 {
				logger.Fmt.Errorf("query target line:%v", err)
			}

			if r_, out_, err_ := cmd.SyncExecBin(
				fmt.Sprintf("lsscsi | awk '{if($1==\"%v\") print $NF}'", uuid)); r_ != 0 || err_ == nil {
				ids = append(ids, ISCSIDriver{
					Host:      strings.TrimSpace(match1[1]),
					Channel:   strings.TrimSpace(match1[2]),
					Target:    strings.TrimSpace(match1[3]),
					Lun:       strings.TrimSpace(match1[4]),
					UUID:      uuid,
					Vendor:    strings.TrimSpace(match2[1]),
					Model:     strings.TrimSpace(match2[2]),
					Rev:       strings.TrimSpace(match2[3]),
					Type:      strings.TrimSpace(match3[1]),
					Revision:  strings.TrimSpace(match3[2]),
					DrivePath: strings.TrimSpace(out_),
					IQNInfo:   iqnInfo,
				})
			}
		}
	}

	return ids, nil
}

/////////////////////////////////////////////////////////////////