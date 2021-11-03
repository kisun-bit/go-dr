package iscsi

import (
	"errors"
	"github.com/kisunSea/go_dr/src/cmd"
	"regexp"
	"strings"
)

var ParseDriverTable = map[string][]string{
	"SessionID":           {"会话 ID", "Session Id"},
	"Initiator":           {"发起程序节点名称", "Initiator Node Name"},
	"TargetNodeName":      {"目标节点名称", "Target Node Name"},
	"TargetName":          {"目标名称", "Target Name"},
	"ISID":                {"ISID", "ISID"},
	"TSID":                {"TSID", "TSID"},
	"ConnNum":             {"连接数量", "Number Connections"},
	"ConnID":              {"连接 ID", "Connection Id"},
	"ConnInitiatorPortal": {"发起程序门户", "Initiator Portal"},
	"ConnTargetPortal":    {"目标门户", "Target Portal"},
	"ConnCID":             {"CID", "CID"},
	"DevTypeDesc":         {"设备类型", "Device Type"},
	"DevNo":               {"设备号", "Device Number"},
	"DevStorageType":      {"存储设备类型", "Storage Device Type"},
	"DevPartNo":           {"分区号", "Partition Number"},
	"DevFriendlyName":     {"友好名称", ""},
	"DevDesc":             {"设备描述", "Device Description"},
	"DevReportMapping":    {"报告的映射", "Reported Mappings"},
	"DevAddr":             {"位置", "Location"},
	"DevInitiatorName":    {"发起程序名称", "Initiator Name"},
	"DevInterName":        {"设备接口名称", "Device Interface Name"},
	"DevRealPath":         {"旧设备名称", "Legacy Device Name"},
	"DevIns":              {"设备实例", "Device Instance"},
}

type ISCSIConn struct {
	ConnID              string // 连接 ID
	ConnInitiatorPortal string // 发起程序门户
	ConnTargetPortal    string // 目标门户
	ConnCID             string // CID
}

type ISCSIDevice struct {
	DevTypeDesc      string // 设备类型
	DevNo            string // 设备号
	DevStorageType   string // 存储设备类型
	DevPartNo        string // 分区号
	DevFriendlyName  string // 友好名称
	DevDesc          string // 磁盘描述
	DevReportMapping string // 报告的映射
	DevAddr          string // 位置
	DevInitiatorName string // 设备接口名称
	DevInterName     string // 目标名称
	DevRealPath      string // 旧设备名称（设备底层真实路径）
	DevIns           string // 设备实例
}

type ISCSIObject struct {
	SessionID      string // 会话ID
	Initiator      string // 发起程序节点名称
	TargetNodeName string // 目标节点名称
	TargetName     string // 目标名称
	ISID           string // ISID
	TSID           string // TSID
	ConnNum        string // 连接数量
	ISCSIConns     []ISCSIConn
	ISCSIDevices   []ISCSIDevice
}

type ISCSIObjs struct {
	LangIdx int
	Objs    []ISCSIObject
}

func QueryISCSIObjs() (is *ISCSIObjs, err error) {
	is = new(ISCSIObjs)
	is.LangIdx = 0 // TODO 暂不支持英文版Windows

	var isInfo string
	if r, out, err := cmd.SyncExecBin("iscsicli sessionlist"); r == 0 && err == nil {
		isInfo = out
	} else {
		return nil, errors.New("failed to query QueryISCSIObjs")
	}

	if strings.TrimSpace(isInfo) == "" {
		return nil, errors.New("iSCSI info is null")
	}

	// 格式化
	ls := strings.Split(isInfo, "\n")
	fixedLines := make([]string, 0)
	for i := 0; i < len(ls); i++ {
		if ls[i] = strings.TrimSpace(ls[i]); ls[i] != "" {
			__ss := strings.Split(ls[i], ":")
			if len(__ss) > 1 {
				__ss[0] = strings.TrimSpace(__ss[0])
				ls[i] = strings.Join(__ss, ":")
			}
			fixedLines = append(fixedLines, ls[i])
		}
	}

	// 统计会话数量
	matchedSessionsNumPattern := regexp.MustCompile("共 (.*)? 次会话")
	msn := matchedSessionsNumPattern.FindStringSubmatch(strings.Join(fixedLines, "\n"))
	if len(msn) < 2 {
		return
	}

	// 组装ISCSIObjs
	var tmpObj *ISCSIObject
	for i := 0; i < len(fixedLines); i++ {
		tmpLine := fixedLines[i]

		if strings.HasPrefix(tmpLine, ParseDriverTable["SessionID"][is.LangIdx]) {
			if tmpObj != nil {
				is.Objs = append(is.Objs, *tmpObj)
				tmpObj = nil
			}
			tmpObj = new(ISCSIObject)
			tmpObj.SessionID = strings.TrimSpace(tmpLine[len(ParseDriverTable["SessionID"][is.LangIdx])+1:])
			continue
		}

		if i == len(fixedLines)-1 {
			if tmpObj != nil {
				is.Objs = append(is.Objs, *tmpObj)
				tmpObj = nil
			}
		}

		if tmpObj == nil {
			continue
		}

		if strings.HasPrefix(tmpLine, ParseDriverTable["Initiator"][is.LangIdx]) {
			tmpObj.Initiator = strings.TrimSpace(tmpLine[len(ParseDriverTable["Initiator"][is.LangIdx])+1:])
		}
		if strings.HasPrefix(tmpLine, ParseDriverTable["TargetNodeName"][is.LangIdx]) {
			tmpObj.TargetNodeName = strings.TrimSpace(tmpLine[len(ParseDriverTable["TargetNodeName"][is.LangIdx])+1:])
		}
		if strings.HasPrefix(tmpLine, ParseDriverTable["TargetName"][is.LangIdx]) {
			tmpObj.TargetName = strings.TrimSpace(tmpLine[len(ParseDriverTable["TargetName"][is.LangIdx])+1:])
		}
		if strings.HasPrefix(tmpLine, ParseDriverTable["ISID"][is.LangIdx]) {
			tmpObj.ISID = strings.TrimSpace(tmpLine[len(ParseDriverTable["ISID"][is.LangIdx])+1:])
		}
		if strings.HasPrefix(tmpLine, ParseDriverTable["TSID"][is.LangIdx]) {
			tmpObj.TSID = strings.TrimSpace(tmpLine[len(ParseDriverTable["TSID"][is.LangIdx])+1:])
		}
		if strings.HasPrefix(tmpLine, ParseDriverTable["ConnNum"][is.LangIdx]) {
			tmpObj.ConnNum = strings.TrimSpace(tmpLine[len(ParseDriverTable["ConnNum"][is.LangIdx])+1:])
		}

		if tmpLine == "连接:" { // 可能存在多个连接
			var tmpConn *ISCSIConn
			tmpConns := make([]ISCSIConn, 0)
			for j := i; j < len(fixedLines); j++ {
				tmpConnLine := fixedLines[j]
				if strings.HasPrefix(tmpConnLine, ParseDriverTable["ConnID"][is.LangIdx]) {
					if tmpConn != nil {
						tmpConns = append(tmpConns, *tmpConn)
						tmpConn = nil
					}
					tmpConn = new(ISCSIConn)
					tmpConn.ConnID = strings.TrimSpace(tmpConnLine[len(ParseDriverTable["ConnID"][is.LangIdx])+1:])
					continue
				}

				if tmpConnLine == "设备:" ||
					strings.HasPrefix(tmpConnLine, ParseDriverTable["SessionID"][is.LangIdx]) ||
					j == len(fixedLines)-1 {
					if tmpConn != nil {
						tmpConns = append(tmpConns, *tmpConn)
						tmpConn = nil
						break
					}
				}

				if tmpConn == nil {
					continue
				}

				if strings.HasPrefix(tmpConnLine, ParseDriverTable["ConnID"][is.LangIdx]) {
					tmpConn.ConnID = strings.TrimSpace(
						tmpConnLine[len(ParseDriverTable["ConnID"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpConnLine, ParseDriverTable["ConnInitiatorPortal"][is.LangIdx]) {
					tmpConn.ConnInitiatorPortal = strings.TrimSpace(
						tmpConnLine[len(ParseDriverTable["ConnInitiatorPortal"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpConnLine, ParseDriverTable["ConnTargetPortal"][is.LangIdx]) {
					tmpConn.ConnTargetPortal = strings.TrimSpace(
						tmpConnLine[len(ParseDriverTable["ConnTargetPortal"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpConnLine, ParseDriverTable["ConnCID"][is.LangIdx]) {
					tmpConn.ConnCID = strings.TrimSpace(
						tmpConnLine[len(ParseDriverTable["ConnCID"][is.LangIdx])+1:])
				}
			}
			if tmpObj != nil {
				tmpObj.ISCSIConns = tmpConns
			}
		}

		if tmpLine == "设备:" { // 可能存在多个设备
			var tmpDev *ISCSIDevice
			tmpDevs := make([]ISCSIDevice, 0)
			for k := i; k < len(fixedLines); k++ {
				tmpDevLine := fixedLines[k]
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevTypeDesc"][is.LangIdx]) {
					if tmpDev != nil {
						tmpDevs = append(tmpDevs, *tmpDev)
						tmpDev = nil
					}
					tmpDev = new(ISCSIDevice)
					tmpDev.DevTypeDesc = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevTypeDesc"][is.LangIdx])+1:])
					continue
				}

				if tmpDevLine == "连接:" || k == len(fixedLines)-1 {
					if tmpDev != nil {
						tmpDevs = append(tmpDevs, *tmpDev)
						tmpDev = nil
						break
					}
				}

				if tmpDev == nil {
					continue
				}

				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevTypeDesc"][is.LangIdx]) {
					tmpDev.DevTypeDesc = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevTypeDesc"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevNo"][is.LangIdx]) {
					tmpDev.DevNo = strings.TrimSpace(tmpDevLine[len(ParseDriverTable["DevNo"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevStorageType"][is.LangIdx]) {
					tmpDev.DevStorageType = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevStorageType"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevPartNo"][is.LangIdx]) {
					tmpDev.DevPartNo = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevPartNo"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevFriendlyName"][is.LangIdx]) {
					tmpDev.DevFriendlyName = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevFriendlyName"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevDesc"][is.LangIdx]) {
					tmpDev.DevDesc = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevDesc"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevReportMapping"][is.LangIdx]) {
					tmpDev.DevReportMapping = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevReportMapping"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevAddr"][is.LangIdx]) {
					tmpDev.DevAddr = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevAddr"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevInitiatorName"][is.LangIdx]) {
					tmpDev.DevInitiatorName = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevInitiatorName"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevInterName"][is.LangIdx]) {
					tmpDev.DevInterName = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevInterName"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevRealPath"][is.LangIdx]) {
					tmpDev.DevRealPath = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevRealPath"][is.LangIdx])+1:])
				}
				if strings.HasPrefix(tmpDevLine, ParseDriverTable["DevIns"][is.LangIdx]) {
					tmpDev.DevIns = strings.TrimSpace(
						tmpDevLine[len(ParseDriverTable["DevIns"][is.LangIdx])+1:])
				}
			}
			if tmpObj != nil {
				tmpObj.ISCSIDevices = tmpDevs
			}
		}
	}

	return is, nil
}

func (i *ISCSIObjs) QueryAttr(block, iqn, attr string) string {
	logger.Fmt.Debugf("block: %v | IQN: %v | attr: %v\n", block, iqn, attr)

	for j := 0; j < len(i.Objs); j++ {
		if i.Objs[j].TargetName != iqn {
			continue
		}

		if block == "" {
			if attr == "SessionID" {
				return i.Objs[j].SessionID
			} else if attr == "DevNo" {
				return i.Objs[j].ISCSIDevices[0].DevNo
			} else if attr == "Initiator" {
				return i.Objs[j].ISCSIDevices[0].DevInitiatorName
			}
		}

		if i.Objs[j].ISCSIDevices != nil && len(i.Objs[j].ISCSIDevices) > 0 {
			for k := 0; k < len(i.Objs[j].ISCSIDevices); k++ {
				if strings.Contains(
					strings.ToLower(i.Objs[j].ISCSIDevices[k].DevFriendlyName),
					strings.ToLower(block)) {
					if attr == "SessionID" {
						return i.Objs[j].SessionID
					} else if attr == "DevNo" {
						return i.Objs[j].ISCSIDevices[k].DevNo
					} else if attr == "Initiator" {
						return i.Objs[j].ISCSIDevices[k].DevInitiatorName
					}
				}
			}
		}
	}
	return ""
}


