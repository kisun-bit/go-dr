package datahandle

import "strconv"

func FmtErrCode2String(code uint32) string {
	return "0x" + strconv.FormatInt(int64(code), 16)
}
