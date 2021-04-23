package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
)

var logger *zap.Logger
func init(){
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	encodeConf := zap.NewProductionEncoderConfig()
	encodeConf.EncodeCaller = zapcore.ShortCallerEncoder
	encodeConf.EncodeTime = zapcore.RFC3339TimeEncoder
	encodeConf.EncodeLevel = zapcore.LowercaseLevelEncoder
}

func Logger() *zap.Logger {
	return logger
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename + ".%Y%m%d%H%M",
		rotatelogs.WithLinkName(filename),
		rotatelogs.WithMaxAge(24 * time.Hour),
		rotatelogs.WithRotationTime(time.Hour))
	if err != nil {
		return nil
	}
	return hook
}
