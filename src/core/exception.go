package core

import (
	"fmt"
	"github.com/kisun-bit/go_dr/src/log"
	"runtime"
	"time"

	"github.com/kisun-bit/go_dr/src/datahandle"
	"github.com/kisun-bit/go_dr/src/meta"
)

// 用于及“标准错误类型”下的“调用栈”
type ExcCtx struct {
	msg      string
	datetime time.Time

	funcName string // 调用方法名
	pkgName  string // 文件路径
	lineNo   int    // 执行时行号
}

// 标准错误类型
type StandardError struct {
	ErrorCode        uint32   // 错误码， 0为“未知错误”类型
	ErrorType        string   // 错误类型
	ErrorDescription string   // 错误描述      用户可读
	ErrorDebug       string   // 错误调试信息  工程师可读
	Trace            []ExcCtx // 调用栈
}

func (jse *StandardError) createErrCtx(msg string, frameLevel int) []ExcCtx {
	_ = jse
	jec := make([]ExcCtx, 0)
	now := time.Now()
	for i := 0; i < frameLevel; i++ {
		pc, file, lineNo, ok := runtime.Caller(i)
		if !ok {
			break
		}

		pcName := runtime.FuncForPC(pc).Name()
		jec = append(
			jec,
			ExcCtx{msg: msg, datetime: now, funcName: pcName, pkgName: file, lineNo: lineNo})
	}
	return jec
}

// @title :  Error 错误信息描述
// @remark:  该方法仅返回错误信息的简短描述，格式：仅包含ErrorDebug及ErrorCode
func (jse *StandardError) Error() string {
	return fmt.Sprintf(
		"！！！！！StandardError ---> ErrCode: \"%s\", \tErrDebug: \"%s\"",
		jse.FmtErrCode2String(), jse.ErrorDebug)
}

// @title :  ErrorDetail 输出标准错误的详细信息
// @remark:  错误信息的详细描述,
//           包含ErrorCode、ErrorDescription、ErrorDebug及Trace
func (jse *StandardError) ErrorDetail() string {
	tLines := "\n----------\nTraceback (most recent call last):\n"
	lastType, lineDots := "", "\t......\n"
	for i, ec := range jse.Trace {
		if lastType != ec.msg {
			lastType = ec.msg
			if i != 0 {
				tLines += lineDots
			}
			tLines += "\t" + lastType + "\n"
		}
		tLines += fmt.Sprintf("\tFile \"<%s>\", Line %d ----> %s\n", ec.pkgName, ec.lineNo, ec.funcName)

	}
	if len(tLines) == 0 {
		tLines += "No Traceback？！！！\n"
	} else {
		tLines += lineDots
	}
	tLines += fmt.Sprintf("####Code: \"%s\", Type: \"%s\"\n####Debug: \"%s\", Desc: \"%s\"\n----------\n",
		jse.FmtErrCode2String(), jse.ErrorType, jse.ErrorDebug, jse.ErrorDescription)
	return tLines
}

func (jse *StandardError) AddMoreDebug(debugStr string) {
	jse.Trace = append(jse.Trace, jse.createErrCtx("[AddMoreDebug] \""+debugStr+"\"", 2)...)
	jse.ErrorDebug += "\t|\t" + debugStr
}

func (jse *StandardError) ChangeDescription(description string) {
	jse.Trace = append(jse.Trace, jse.createErrCtx(
		fmt.Sprintf("[ChangeDescription] cvt \"%s\" to \"%s\"", jse.ErrorDescription, description),
		2)...)
	jse.ErrorDescription = description
}

func (jse *StandardError) FmtErrCode2String() string {
	return datahandle.FmtErrCode2String(jse.ErrorCode)
}

// @title :  RaiseStandardError 生成定制化的"标准错误"
// @params:  code    错误码
//           errType 错误类型
//	         reason  错误原因，用户阅读
//	         debug   错误调试信息
// @return:  *StandardError
func RaiseStandardError(code uint32, errType, reason, debug string) (err *StandardError) {
	jse := StandardError{
		ErrorType:        errType,
		ErrorCode:        code,
		ErrorDescription: reason,
		ErrorDebug:       debug,
		Trace:            make([]ExcCtx, 0),
	}

	jse.Trace = append(jse.Trace, jse.createErrCtx(
		fmt.Sprintf("[Tance] %v|%s", jse.FmtErrCode2String(), jse.ErrorType), meta.ErrDefaultFrameLevel)...)
	return &jse
}

// @title :  StandardizeErr 将错误转换为标准错误
func StandardizeErr(err error) (jse *StandardError) {
	return RaiseStandardError(meta.ErrInternal, "InternalError",
		"内部错误，错误代码："+datahandle.FmtErrCode2String(meta.ErrInternal), err.Error(),
	)
}

// @title : ConvertPanic2StandardErr 将已捕获的panic异常装换为标准错误
func ConvertPanic2StandardErr(r interface{}) *StandardError {
	var jse_ *StandardError

	switch r.(type) {
	case *StandardError:
		jse_ = r.(*StandardError)
	case StandardError:
		jseTmp := r.(StandardError)
		jse_ = &jseTmp

	// TODO 更多错误类型 gRPC相关、第三方组件错误类型

	default:
		jse_ = RaiseStandardError(
			meta.ErrInternal,
			"PanicError",
			"内部异常， 错误代码："+datahandle.FmtErrCode2String(meta.ErrInternal),
			fmt.Sprintf("%v\r\n", r))
	}
	return jse_
}

// @title :  CatchPanicErr 捕获panic异常并将其转化为标准错误
// @remark:  该方法使用在`defer`中
// @exp   :  defer CatchPanicErr(logger)
func CatchPanicErr() *StandardError {
	err := recover()
	if err == nil {
		log.DLogger.Error("err is nil？！！！")
		return nil
	}

	jse_ := ConvertPanic2StandardErr(err)
	log.DLogger.Warn(jse_.ErrorDetail())
	return jse_
}

func StandardPanic(code uint32, errType, reason, debug string) {
	jse := RaiseStandardError(code, errType, reason, debug)
	log.DLogger.Error(jse.ErrorDetail())
	panic(jse)
}
