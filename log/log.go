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

var loggerZap *zap.Logger
func Init(lvl string, path string){
	if loggerZap == nil {
		loggerZap, _ = zap.NewProduction()
	}
	encodeConf := zap.NewProductionEncoderConfig()
	encodeConf.EncodeCaller = zapcore.ShortCallerEncoder
	encodeConf.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encodeConf.EncodeLevel = zapcore.LowercaseLevelEncoder
	encodeConf.TimeKey = "time"
	errorLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.ErrorLevel
	})

	debugLevel := zap.LevelEnablerFunc(func(level zapcore.Level) bool {
		return level >= zap.DebugLevel
	})

	errorWritter := getWriter(path + "/error")
	debugWritter := getWriter(path + "/debug")

	core := zapcore.NewTee(
			zapcore.NewCore(zapcore.NewConsoleEncoder(encodeConf), zapcore.NewMultiWriteSyncer(zapcore.AddSync(errorWritter), zapcore.AddSync(os.Stdout)), errorLevel),
			zapcore.NewCore(zapcore.NewConsoleEncoder(encodeConf), zapcore.NewMultiWriteSyncer(zapcore.AddSync(debugWritter), zapcore.AddSync(os.Stdout)), debugLevel),
		)
	loggerZap = zap.New(core, zap.AddCaller())

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

func Logger() *zap.SugaredLogger {
	return loggerZap.Sugar()
}

func getWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename + ".%Y%m%d%H.log",
		rotatelogs.WithLinkName(filename + ".log"),
		rotatelogs.WithMaxAge(24 * time.Hour),
		rotatelogs.WithRotationTime(time.Hour))
	if err != nil {
		return nil
	}
	return hook
}
