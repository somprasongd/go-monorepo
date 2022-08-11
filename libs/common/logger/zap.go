package logger

import "go.uber.org/zap"

type zapLogger struct {
	*zap.Logger
}

func newZapLogger(logger *zap.Logger) Logger {
	return &zapLogger{logger}
}

func (l *zapLogger) Info(msg string) {
	l.Logger.Info(msg)
}

func (l *zapLogger) InfoWithFields(msg string, fileds map[string]interface{}) {
	zf := l.toZapFields(fileds)
	l.Logger.Info(msg, zf...)
}

func (l *zapLogger) Debug(msg string) {
	l.Logger.Debug(msg)
}

func (l *zapLogger) DebugWithFields(msg string, fileds map[string]interface{}) {
	zf := l.toZapFields(fileds)
	l.Logger.Debug(msg, zf...)
}

func (l *zapLogger) Error(msg string) {
	l.Logger.Error(msg)
}

func (l *zapLogger) ErrorWithFields(msg string, fileds map[string]interface{}) {
	zf := l.toZapFields(fileds)
	l.Logger.Error(msg, zf...)
}

func (l *zapLogger) Warn(msg string) {
	l.Logger.Warn(msg)
}

func (l *zapLogger) WarnWithFields(msg string, fileds map[string]interface{}) {
	zf := l.toZapFields(fileds)
	l.Logger.Warn(msg, zf...)
}

func (l *zapLogger) Fatal(msg string) {
	l.Logger.Fatal(msg)
}

func (l *zapLogger) FatalWithFields(msg string, fileds map[string]interface{}) {
	zf := l.toZapFields(fileds)
	l.Logger.Fatal(msg, zf...)
}

func (l *zapLogger) Panic(msg string) {
	l.Logger.Panic(msg)
}

func (l *zapLogger) PanicWithFields(msg string, fileds map[string]interface{}) {
	zf := l.toZapFields(fileds)
	l.Logger.Panic(msg, zf...)
}

func (l *zapLogger) toZapFields(fileds map[string]interface{}) []zap.Field {
	fz := []zap.Field{}
	for k, v := range fileds {
		fz = append(fz, zap.Any(k, v))
	}
	return fz
}
