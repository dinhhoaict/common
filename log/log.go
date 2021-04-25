package log

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"strings"
	"time"
)

const (
	DEBUG = "debug"
	ERROR = "error"
	INFO = "info"
)

var logger *zap.Logger
func Init(lvl string, path string){
	if logger == nil {
		logger, _ = zap.NewProduction()
	}
	encodeConf := zap.NewProductionEncoderConfig()
	encodeConf.EncodeCaller = zapcore.ShortCallerEncoder
	encodeConf.EncodeTime = zapcore.RFC3339TimeEncoder
	encodeConf.EncodeLevel = zapcore.LowercaseLevelEncoder
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})

	debugLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.DebugLevel
	})

	errorWritter := getWriter(path + "./error.log")
	debugWritter := getWriter(path + "./debug.log")

	core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewJSONEncoder(encodeConf), zapcore.NewMultiWriteSyncer(zapcore.AddSync(errorWritter), zapcore.AddSync(os.Stdout)), errorLevel),
			zapcore.NewCore(zapcore.NewJSONEncoder(encodeConf), zapcore.NewMultiWriteSyncer(zapcore.AddSync(debugWritter), zapcore.AddSync(os.Stdout)), debugLevel),
		)
	logger = zap.New(core, zap.AddCaller())

}

func getLevel(lvl string) zapcore.Level {
	switch strings.ToLower(lvl) {
	case ERROR:
		return zapcore.ErrorLevel
	case DEBUG:
		return zapcore.DebugLevel
	case INFO:
		return zapcore.InfoLevel
	}
	return zapcore.DebugLevel
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
