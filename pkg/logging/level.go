package logging

import "go.uber.org/zap/zapcore"

type LogLevel string

const (
	levelDebug   LogLevel = "DEBUG"
	levelInfo    LogLevel = "INFO"
	levelWarning LogLevel = "WARNING"
	levelError   LogLevel = "ERROR"
)

func levelToZapLevel(level LogLevel) zapcore.Level {
	switch level {
	case levelDebug:
		return zapcore.DebugLevel
	case levelInfo:
		return zapcore.InfoLevel
	case levelWarning:
		return zapcore.WarnLevel
	case levelError:
		return zapcore.ErrorLevel
	}

	return zapcore.DebugLevel
}

func levelEncoder() zapcore.LevelEncoder {
	return func(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		switch l {
		case zapcore.DebugLevel:
			enc.AppendString(string(levelDebug))
		case zapcore.InfoLevel:
			enc.AppendString(string(levelInfo))
		case zapcore.WarnLevel:
			enc.AppendString(string(levelWarning))
		case zapcore.ErrorLevel:
			enc.AppendString(string(levelError))
		}
	}
}
