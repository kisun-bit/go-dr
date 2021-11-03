package log

import (
	"go.uber.org/zap/zapcore"
)

const _defaultTime = "2006/01/02 - 15:04:05.000"

// Logger Level
type Level = zapcore.Level

const (
	LDebug      = zapcore.DebugLevel
	LInfo       = zapcore.InfoLevel
	LWarn       = zapcore.WarnLevel
	LError      = zapcore.ErrorLevel
	LDebugPanic = zapcore.DPanicLevel
	LPanic      = zapcore.PanicLevel
	LFatal      = zapcore.FatalLevel
)

// Logger Encoder
type Encoder string

var (
	EJson    Encoder = "json"
	EConsole Encoder = "console"
)

// Logger Type
type lType string

var (
	lTypeRateFile lType = "RateFileStream"
	lTypeConsole  lType = "ConsoleStream"
)

// *********************************************
// create a logger ... Rate File or Console
// *********************************************

// NewRateFileLimitAgeSugaredLogger return a custom rate file logger (limit the age of log files)
func NewRateFileLimitAgeSugaredLogger(_name, _path string, _level Level, _encoder Encoder, _rotationTime, _retentionTime int) (
	log *Logger, err error) {

	// init global Logger
	defer func() {
		DLogger = log
	}()

	return New(Config{
		lType_:          lTypeRateFile,
		Color:           false,
		LogFilePath:     _path,
		MaxRotationTime: _rotationTime,
		MaxAge:          _retentionTime,
		TimeFormat:      _defaultTime,
		Stacktrace:      LError,
		ShowLine:        true,
		Encoder:         _encoder,
		Level:           _level,
		Prefix:          _name,
	})
}

// TODO abandon ...
//// NewRateFileLimitSizeLogger return a custom rate file logger (limit the size of log files)
//func NewRateFileLimitSizeLogger(_path string, _encoder Encoder, _filesMaxSize int64, _keepFilesCount int) (
//	log *SugarLogger, err error) {
//
//	return New(Config{
//		lType_:           lTypeRateFile,
//		Color:            false,
//		LogFilePath:      _path,
//		MaxSize:          _filesMaxSize,
//		MaxRotationCount: _keepFilesCount,
//		TimeFormat:       _defaultTime,
//		Stacktrace:       LError,
//		ShowLine:         true,
//		Encoder:          _encoder,
//	})
//}

// NewConsoleLogger return a custom console logger
func NewConsoleLogger(_prefix string, _level Level, _encoder Encoder) (log *Logger, err error) {
	return New(Config{
		lType_:     lTypeConsole,
		Color:      true,
		TimeFormat: _defaultTime,
		Stacktrace: LError,
		ShowLine:   true,
		Encoder:    _encoder,
		Prefix:     _prefix,
		Level:      _level,
	})
}

// NewRateFileLoggerByConf TODO ... return a custom rate file logger by config file

// NewConsoleLoggerByConf TODO ... return a custom console logger by config file
