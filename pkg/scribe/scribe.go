package scribe

import (
	"cyclic/pkg/colonel"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

var Scribe *zap.Logger

// Level is a map of string to zapcore.Level
var Level = map[string]zapcore.Level{
	"debug":  zapcore.DebugLevel,
	"info":   zapcore.InfoLevel,
	"warn":   zapcore.WarnLevel,
	"error":  zapcore.ErrorLevel,
	"dpanic": zapcore.DPanicLevel,
	"panic":  zapcore.PanicLevel,
	"fatal":  zapcore.FatalLevel,
}

// Init initializes the logger Scribe
func Init() {
	Scribe = build(Level[colonel.Writ.Logger.Level],
		colonel.Writ.Logger.File,
		colonel.Writ.Logger.MaxSize,
		colonel.Writ.Logger.MaxBackups,
		colonel.Writ.Logger.MaxAge,
		colonel.Writ.Logger.Compress)
}

// New creates a new zap logger
func build(level zapcore.Level, file string, maxSize int, maxBackups int, maxAge int, compress bool) *zap.Logger {
	fileEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder, // file no need color
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
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // console need color
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	fileEncoder := zapcore.NewJSONEncoder(fileEncoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(consoleEncoderConfig)

	// create zap core
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, zapcore.AddSync(&lumberjack.Logger{
			Filename:   file,
			MaxSize:    maxSize,
			MaxBackups: maxBackups,
			MaxAge:     maxAge,
			Compress:   compress,
		}), level),
		zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), level),
	)

	var options []zap.Option // options for logger

	// if the server mode is debug, add development mode, caller(the file and line number), and stacktrace
	// this is no matter the level of the logger
	if colonel.Writ.Server.Mode == "debug" {
		options = append(options, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// create logger
	logger := zap.New(core, options...)

	return logger
}
