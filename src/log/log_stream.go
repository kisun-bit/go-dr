package log

import (
	"fmt"
	"os"
	"time"

	rotateLogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger ...
type Logger struct {
	*zap.Logger

	format  Encoder // encoder format
	console bool
	config  zapcore.EncoderConfig
	level   zapcore.Level
	writer  zapcore.WriteSyncer
	Fmt     *zap.SugaredLogger
}

// Config ...
type Config struct {
	lType_      lType   // logger type , "rate file" or "console"
	LogFilePath string  // prefix on the name of rate log file
	Level       Level   // the lowest level of log printing , default "info"
	Prefix      string  // prefix on each line to identify the logger
	TimeFormat  string  // the time format of each line in the log
	Color       bool    // the color of log level
	ShowLine    bool    // show log call line number
	Stacktrace  Level   // stack trace log level
	Encoder     Encoder // log encoding format, divided into "json" and "console", default "console"

	// Both time-limiting and size-limiting
	// are the ways to limit the log file size to unlimited growth.
	// Only one of these two methods can be used as your mode, not both.

	// mode 1: time limiting
	MaxAge          int // maximum retention time , default "1" hour
	MaxRotationTime int // file cutting time , default "1"

	// mode 2: size limiting TODO abandon...  As of now, only mode 1 is supported
	MaxSize          int64 // single log file's maximum retention size
	MaxRotationCount int   // maximum number of retained copies
}

// New ...
func New(c Config) (log *Logger, err error) {
	log = &Logger{
		format:  c.Encoder,
		console: c.lType_ == lTypeConsole,
		level:   convertLevel(c.Level),
	}
	if c.Prefix != "" {
		c.Prefix += " "
	}
	if c.TimeFormat == "" {
		c.TimeFormat = _defaultTime
	}
	log.config = initEncoderConfig(c.Prefix + c.TimeFormat)
	if c.Stacktrace == Level(-1) {
		log.config.StacktraceKey = ""
	}
	if !c.Color {
		// close color
		log.config.EncodeLevel = zapcore.CapitalLevelEncoder
	}

	if err = log.setWriter(c.LogFilePath, c.MaxAge, c.MaxRotationTime); err != nil {
		fmt.Printf("set writer failed: %v", err.Error())
		return
	}

	log.Logger = zap.New(log.getEncoderCore(log.config), zap.AddStacktrace(convertLevel(c.Stacktrace)))
	if c.ShowLine {
		log.Logger = log.WithOptions(zap.AddCaller())
	}

	log.Fmt = log.Logger.Sugar()
	return
}

// GinLogConfig return gin web framework log configuration
//func (l *Logger) GinLogConfig() gin.LoggerConfig {
//	return gin.LoggerConfig{
//		Output:    NewGinLogger(l.Logger),
//		Formatter: GinFormatter,
//	}
//}

// Showline configures the SugarLogger to annotate each message with the filename, line number, and function name of zap's caller.
//func (l *SugarLogger) NoShowline() *zap.SugarLogger {
//	c := l.config
//	c. = zapcore.CapitalLevelEncoder
//	return l.wrapCore(c)
//}

// SetStacktraceLevel configures the SugarLogger to record a stack trace for all messages at or above a given level.
func (l *Logger) SetStacktraceLevel(level Level) *zap.Logger {
	return l.WithOptions(zap.AddStacktrace(convertLevel(level)))
}

// CloseStacktrace ...
func (l *Logger) CloseStacktrace() *zap.Logger {
	c := l.config
	c.StacktraceKey = ""
	return l.wrapCore(c)
}

// SetTimeFormat sets the log output format.
// default time format is `2006/01/02 - 15:04:05.000`,
func (l *Logger) SetTimeFormat(timeFormat string) *zap.Logger {
	c := l.config
	c.EncodeTime = customTimeEncoder(timeFormat)
	return l.wrapCore(c)
}

// NoColor ...
func (l *Logger) NoColor() *zap.Logger {
	c := l.config
	c.EncodeLevel = zapcore.CapitalLevelEncoder
	return l.wrapCore(c)
}

// SetJSONStyle change output style to json
func (l *Logger) SetJSONStyle() *zap.Logger {
	tmp := l.format
	defer func() { l.format = tmp }()
	l.format = EJson
	return l.wrapCore(l.config)
}

// wrapCore wraps or replaces the SugarLogger's underlying zapcore.Core.
func (l *Logger) wrapCore(ec zapcore.EncoderConfig) *zap.Logger {
	return l.WithOptions(zap.WrapCore(func(c zapcore.Core) zapcore.Core {
		return l.getEncoderCore(ec)
	}))
}

// setWriter zap logger's writer use file-rotateLogs
func (l *Logger) setWriter(filePath_ string, maxAge, maxRotationTime int) error {
	if l.console {
		l.writer = zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout))
		return nil
	}

	var (
		fileWriter *rotateLogs.RotateLogs
		err        error
	)

	if maxAge != 0 && maxRotationTime != 0 { // limited by age
		fileWriter, err = rotateLogs.New(
			filePath_+"_%Y%m%d%H%M.log",
			rotateLogs.WithLinkName(filePath_),
			rotateLogs.WithMaxAge(time.Duration(maxAge)*24*time.Hour),
			rotateLogs.WithRotationTime(time.Duration(maxRotationTime)*24*time.Hour),
		)
	}

	if err != nil {
		return err
	}
	l.writer = zapcore.AddSync(fileWriter)
	return nil
}

// getEncoderCore uses the new config to get the new core
func (l *Logger) getEncoderCore(ec zapcore.EncoderConfig) (core zapcore.Core) {
	var encoder zapcore.Encoder
	if l.format == EJson {
		encoder = zapcore.NewJSONEncoder(ec)
	} else if l.format == EConsole {
		encoder = zapcore.NewConsoleEncoder(ec)
	} else {
		encoder = zapcore.NewConsoleEncoder(ec)
	}
	return zapcore.NewCore(encoder, l.writer, l.level)
}

// initEncoderConfig init zapcore.EncoderConfig
func initEncoderConfig(format string) zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     customTimeEncoder(format),
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
}

// customTimeEncoder sets custom log output time format
func customTimeEncoder(format string) zapcore.TimeEncoder {
	return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format(format))
	}
}

func convertLevel(lvl Level) zapcore.Level {
	return lvl
}
