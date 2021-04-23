package log

import (
	"strconv"
	"time"

	"github.com/natefinch/lumberjack"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/kisunSea/jpkt/src/core"
)

type LgConf struct {
	filename     string
	maxSize      int
	maxBackups   int
	maxAge       int
	compress     bool
	level        string
	splitByLevel bool
}

type JLogger struct {
	logConf LgConf
	logger  *zap.SugaredLogger
}

func GetJLoggerByMapConf(logConf map[string]interface{}) *zap.SugaredLogger {
	jl := &JLogger{
		logConf: parseLogConf(logConf),
		logger:  nil,
	}
	jl.setConf()
	return jl.logger
}

func GetJLoggerByConf(baseDir, confFileName, LoggerName string) *zap.SugaredLogger {
	jl := &JLogger{
		logConf: readLogConf(baseDir, confFileName, LoggerName),
		logger:  nil,
	}
	jl.setConf()
	return jl.logger
}

func (jl *JLogger) setConf() {
	writeSyncer := jl.getLogWriter()
	encoder := jl.getEncoder()
	_core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	jl.logger = zap.New(_core, zap.AddCaller()).Sugar()
}

func parseLogConf(lc map[string]interface{}) LgConf {
	return LgConf{
		filename:     lc["Filename"].(string),
		maxSize:      lc["MaxSize"].(int),
		maxBackups:   lc["MaxBackups"].(int),
		maxAge:       lc["MaxAge"].(int),
		compress:     lc["Compress"].(bool),
		level:        lc["Level"].(string),
		splitByLevel: lc["SplitByLevel"].(bool),
	}
}

/**/
func readLogConf(baseDir, confFileName, LoggerName string) LgConf {
	v := viper.New()
	v.SetConfigName(confFileName)
	v.AddConfigPath(baseDir)
	err := v.ReadInConfig()
	if err != nil {
		panic(err)
	}

	var __LoggerItem = func(groupName string) func(itemName string) string {
		return func(itemName string) string {
			return groupName + "." + itemName
		}
	}

	__getItem := __LoggerItem(LoggerName)
	logConf := LgConf{
		filename: v.GetString(
			__getItem("Filename")),
		maxSize: v.GetInt(
			__getItem("MaxSize")),
		maxBackups: v.GetInt(
			__getItem("MaxBackups")),
		maxAge: v.GetInt(
			__getItem("MaxAge")),
		compress: v.GetBool(
			__getItem("Compress")),
		level: v.GetString(
			__getItem("Level")),
		splitByLevel: v.GetBool(
			__getItem("SplitByLevel")),
	}

	return logConf
}

func (jl *JLogger) getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = func(i time.Time, encoder zapcore.PrimitiveArrayEncoder) {
		encoder.AppendString(i.Format("2006-01-02 15:04:05.000") + "\t" + strconv.Itoa(core.Gid()))
	}
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func (jl *JLogger) getLogWriter() zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   jl.logConf.filename,
		MaxSize:    jl.logConf.maxSize,
		MaxBackups: jl.logConf.maxBackups,
		MaxAge:     jl.logConf.maxAge,
		Compress:   jl.logConf.compress,
	}
	return zapcore.AddSync(lumberJackLogger)
}
