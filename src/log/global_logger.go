package log

var DLogger, _ = NewConsoleLogger("GO-DR", LDebug, EConsole)

func SetPacketLogger(l *Logger) {
	DLogger = l
	l.Fmt.Debug("convert go-dr-packet logger is ok")
}
