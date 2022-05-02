// Package log contains a global Zap logger.
package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var SugaredLogger *zap.SugaredLogger

func init() {
	zcfg := zap.NewProductionConfig()

	zcfg.Level.SetLevel(zap.InfoLevel)
	zcfg.Encoding = "console"
	zcfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	zcfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := zcfg.Build()
	if err != nil {
		panic(err)
	}

	zap.RedirectStdLog(logger)

	Logger = logger
	SugaredLogger = Logger.WithOptions(zap.AddCallerSkip(1)).Sugar()
}

func Debug(v ...any) {
	SugaredLogger.Debug(v...)
}

func Info(v ...any) {
	SugaredLogger.Info(v...)
}

func Warn(v ...any) {
	SugaredLogger.Warn(v...)
}

func Error(v ...any) {
	SugaredLogger.Error(v...)
}

func Fatal(v ...any) {
	SugaredLogger.Fatal(v...)
}

func Debugf(tmpl string, v ...any) {
	SugaredLogger.Debugf(tmpl, v...)
}

func Infof(tmpl string, v ...any) {
	SugaredLogger.Infof(tmpl, v...)
}

func Warnf(tmpl string, v ...any) {
	SugaredLogger.Warnf(tmpl, v...)
}

func Errorf(tmpl string, v ...any) {
	SugaredLogger.Errorf(tmpl, v...)
}

func Fatalf(tmpl string, v ...any) {
	SugaredLogger.Fatalf(tmpl, v...)
}
