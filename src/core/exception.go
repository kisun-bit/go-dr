package core

import (
	"fmt"
	"go.uber.org/zap"
	"jpkt/src/datahandle"
	"jpkt/src/meta"
	"runtime"
	"time"
)

// 用于及“标准错误类型”下的“调用栈”
type JpktExcCtx struct {
	msg      string
	datetime time.Time

	funcName string // 调用方法名
	pkgName  string // 文件路径
	lineNo   int    // 执行时行号
}

// 标准错误类型
type JpktStandardError struct {
	ErrorCode        uint32       // 错误码， 0为“未知错误”类型
	ErrorType        string       // 错误类型
	ErrorDescription string       // 错误描述      用户可读
	ErrorDebug       string       // 错误调试信息  工程师可读
	Trace            []JpktExcCtx // 调用栈
}

func (jse *JpktStandardError) createErrCtx(msg string, frameLevel int) []JpktExcCtx {
	_ = jse
	jec := make([]JpktExcCtx, 0)
	now := time.Now()
	for i := 0; i < frameLevel; i++ {
		pc, file, lineNo, ok := runtime.Caller(i)
		if !ok {
			break
		}

		pcName := runtime.FuncForPC(pc).Name()
		jec = append(
			jec,
			JpktExcCtx{msg: msg, datetime: now, funcName: pcName, pkgName: file, lineNo: lineNo,})
	}
	return jec
}

// @title :  Error 错误信息描述
// @remark:  该方法仅返回错误信息的简短描述，格式：仅包含ErrorDebug及ErrorCode
func (jse *JpktStandardError) Error() string {
	return fmt.Sprintf(
		"！！！！！JpktStandardError ---> ErrCode: \"%s\", \tErrDebug: \"%s\"",
		jse.FmtErrCode2String(), jse.ErrorDebug)
}

// @title :  ErrorDetail 输出标准错误的详细信息
// @remark:  错误信息的详细描述,
//           包含ErrorCode、ErrorDescription、ErrorDebug及Trace
func (jse *JpktStandardError) ErrorDetail() string {
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

func (jse *JpktStandardError) AddMoreDebug(debugStr string) {
	jse.Trace = append(jse.Trace, jse.createErrCtx("[AddMoreDebug] \""+debugStr+"\"", 2)...)
	jse.ErrorDebug += "\t|\t" + debugStr
}

func (jse *JpktStandardError) ChangeDescription(description string) {
	jse.Trace = append(jse.Trace, jse.createErrCtx(
		fmt.Sprintf("[ChangeDescription] cvt \"%s\" to \"%s\"", jse.ErrorDescription, description),
		2)...)
	jse.ErrorDescription = description
}

func (jse *JpktStandardError) FmtErrCode2String() string {
	return datahandle.FmtErrCode2String(jse.ErrorCode)
}

// @title :  RaiseStandardError 生成定制化的"标准错误"
// @params:  code    错误码
//           errType 错误类型
//	         reason  错误原因，用户阅读
//	         debug   错误调试信息
// @return:  *JpktStandardError
func RaiseStandardError(code uint32, errType, reason, debug string) (err *JpktStandardError) {
	jse := JpktStandardError{
		ErrorType:        errType,
		ErrorCode:        code,
		ErrorDescription: reason,
		ErrorDebug:       debug,
		Trace:            make([]JpktExcCtx, 0),
	}

	jse.Trace = append(jse.Trace, jse.createErrCtx(
		fmt.Sprintf("[Tance] %v|%s", jse.FmtErrCode2String(), jse.ErrorType), meta.ErrDefaultFrameLevel)...)
	return &jse
}

// @title :  StandardizeErr 将错误转换为标准错误
func StandardizeErr(err error) (jse *JpktStandardError) {
	return RaiseStandardError(meta.JErrInternal, "InternalError",
		"内部错误，错误代码："+datahandle.FmtErrCode2String(meta.JErrInternal), err.Error(),
	)
}

// @title :  CatchPanicErr 捕获panic异常并将其转化为标准错误
// @remark:  该方法使用在`defer`中
// @exp   :  defer CatchPanicErr(logger)
func CatchPanicErr(logger *zap.Logger) *JpktStandardError {
	err := recover()
	if err == nil {
		if logger != nil {
			logger.Error("err is nil？！！！")
		}
		return nil
	}

	var jse_ *JpktStandardError

	switch err.(type) {
	case *JpktStandardError:
		jse_ = err.(*JpktStandardError)
	case JpktStandardError:
		jseTmp := err.(JpktStandardError)
		jse_ = &jseTmp
	// TODO 更多错误类型 gRPC相关、第三方组件错误类型
	default:
		jse_ = RaiseStandardError(
			meta.JErrInternal,
			"PanicError",
			"内部异常， 错误代码："+datahandle.FmtErrCode2String(meta.JErrInternal),
			fmt.Sprintf("%v\r\n", err), )
	}

	if logger != nil {
		logger.Warn(jse_.ErrorDetail())
	}
	return jse_
}

func StandardPanic(code uint32, errType, reason, debug string) {
	jse := RaiseStandardError(code, errType, reason, debug)
	panic(jse)
}
