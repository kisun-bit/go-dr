package core

import (
	"bytes"
	"encoding/json"
	"runtime"
	"strconv"
	"strings"
)

// Gid return the goroutine ID
func Gid() (gid int) {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	idField := strings.Fields(strings.TrimPrefix(string(buf[:n]), "goroutine "))[0]
	gid, err := strconv.Atoi(idField)
	if err != nil {
		gid = -1
	}
	return gid
}

// PrettyPrintStruct prints the visible properties of the structure in JSON format
func PrettyPrintStruct(v interface{}) string {
	b, err := json.Marshal(v)
	if err != nil {
		return "jsonMarshal err:\t" + err.Error()
	}

	var out bytes.Buffer
	err = json.Indent(&out, b, "", "  ")
	if err != nil {
		return "jsonIndent err:\t" + err.Error()
	}

	return out.String()
}
