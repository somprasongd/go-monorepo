package logger

import (
	"os"

	"go.elastic.co/ecszap"
	"go.uber.org/zap"
)

type Logger interface {
	Info(msg string)
	InfoWithFields(msg string, fileds map[string]interface{})
	Debug(msg string)
	DebugWithFields(msg string, fileds map[string]interface{})
	Error(msg string)
	ErrorWithFields(msg string, fileds map[string]interface{})
	Warn(msg string)
	WarnWithFields(msg string, fileds map[string]interface{})
	Fatal(msg string)
	FatalWithFields(msg string, fileds map[string]interface{})
	Panic(msg string)
	PanicWithFields(msg string, fileds map[string]interface{})
}

var zlog *zap.Logger
var Default Logger

func init() {
	var err error

	mode := os.Getenv("APP_MODE")
	mode = "production"
	var config zap.Config
	if mode == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
	}

	config.EncoderConfig = ecszap.ECSCompatibleEncoderConfig(config.EncoderConfig)

	zlog, err = config.Build(ecszap.WrapCoreOption(), zap.AddCallerSkip(1))

	if err != nil {
		panic(err)
	}

	Default = newZapLogger(zlog)
}

func NewWithFields(fileds map[string]interface{}) Logger {
	opts := []zap.Field{}
	for k, v := range fileds {
		opts = append(opts, zap.Any(k, v))
	}

	nlog := zlog.With(opts...)
	return newZapLogger(nlog)
}
