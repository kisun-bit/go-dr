package main

import (
	"errors"
	"github.com/kisun-bit/go-dr/src/log"
	"go.uber.org/zap"
)

func DemoConsoleLogger() {
	log_, err := log.NewConsoleLogger("1111", log.LDebug, log.EConsole)
	if err != nil {
		panic(err)
	}

	log_.Debug("debug...")
	log_.Info("info...")
	log_.Error("error...")
	log_.Fmt.Errorf("=====================: %s", "ssssssss")
	log_.SetStacktraceLevel(log.LDebug).Info("info stack...")
	log_.NoColor().Info("info no color...")
	log_.SetJSONStyle().Warn("warn json...",
		zap.String("start from", "14:00"), zap.String("probability ", "60%"))
	log_.CloseStacktrace().Fatal("fatal close stack trace", zap.Error(errors.New("not found")))
	log_.Error("finish...")
}

func DemoFileLogger() {

	// 最大保留30天的日志，且每隔两天便分割日志
	log_, err := log.NewRateFileLimitAgeSugaredLogger(
		`D:\workspace\go_dr\demo\log\dr.log`, log.LDebug, log.EConsole, 2, 30)
	if err != nil {
		panic(err)
	}

	log_.Debug("debug...dasfdasfffffffffffffffffasfasdffffffffffffffffffffffff")
	log_.Info("info...ffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	log_.Error("error...ffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	log_.Error("finish...fffffffffffffffffffffffffffffffffffffffffffffffffffff")
}

func main() {
	//DemoFileLogger()
	DemoConsoleLogger()
}
