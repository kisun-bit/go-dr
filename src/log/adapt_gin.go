package log

import (
	"bytes"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Writer ...
type Writer struct {
	logFunc func(msg string, fields ...zapcore.Field)
}

// NewGinLogger ...
func NewGinLogger(log *zap.Logger) *Writer {
	logger := log.WithOptions(zap.AddCallerSkip(3))
	return &Writer{logger.Info}
}

func (l *Writer) Write(p []byte) (int, error) {
	p = bytes.TrimSpace(p)
	l.logFunc(string(p))
	return len(p), nil
}

// GinFormatter ...
func GinFormatter(param gin.LogFormatterParams) string {
	gin.ForceConsoleColor()
	var statusColor, methodColor, resetColor string
	if param.IsOutputColor() {
		statusColor = param.StatusCodeColor()
		methodColor = param.MethodColor()
		resetColor = param.ResetColor()
	}
	if param.Latency > time.Minute {
		// Truncate in a golang < 1.8 safe way
		param.Latency = param.Latency - param.Latency%time.Second
	}
	return fmt.Sprintf("%s %3d %s| %13v | %15s |%s %-7s %s %#v\n%s",
		statusColor, param.StatusCode, resetColor,
		param.Latency,
		param.ClientIP,
		methodColor, param.Method, resetColor,
		param.Path,
		param.ErrorMessage,
	)
}
