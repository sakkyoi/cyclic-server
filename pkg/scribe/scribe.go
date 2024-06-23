package scribe

import (
	"cyclic/pkg/colonel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Scribe *zap.Logger

var Level = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

func Init() {
	Scribe = New(Level[colonel.Writ.Logger.Level],
		colonel.Writ.Logger.File,
		colonel.Writ.Logger.MaxSize,
		colonel.Writ.Logger.MaxBackups,
		colonel.Writ.Logger.MaxAge,
		colonel.Writ.Logger.Compress)
}

func New(level zapcore.Level, file string, maxSize int, maxBackups int, maxAge int, compress bool) *zap.Logger {
	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
	consoleEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(&lumberjack.Logger{
			Filename:   file,
			MaxSize:    maxSize, // megabytes
			MaxBackups: maxBackups,
			MaxAge:     maxAge, // days
			Compress:   compress,
		}), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	var options []zap.Option

	if colonel.Writ.Server.Mode == "debug" {
		options = append(options, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	logger := zap.New(core, options...)

	return logger
}
