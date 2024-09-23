package logger

import (
	"io"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var logFile *os.File

type Config struct {
	Encoder       zapcore.Encoder
	LevelEnabler  func(zapcore.Level) bool
	Writer        io.Writer
	RegisterHooks []func(entry zapcore.Entry) error
}

func Init(cfgs ...Config) {
	if Logger != nil {
		Close()
	}

	cores := []zapcore.Core{}
	for _, c := range cfgs {
		var encoder zapcore.Encoder
		var levelEnabler func(zapcore.Level) bool
		var output zapcore.WriteSyncer

		if c.Encoder == nil {
			encoder = zapcore.NewConsoleEncoder(GetDefaultEncoderConfig())
		} else {
			encoder = c.Encoder
		}
		if c.LevelEnabler == nil {
			levelEnabler = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
				/*
					Send all levels to log
				*/
				return true
			})
		} else {
			levelEnabler = c.LevelEnabler
		}
		if c.Writer == nil {
			output = os.Stdout
		} else {
			output = zapcore.AddSync(c.Writer)
		}
		newCore := zapcore.NewCore(
			encoder,
			zapcore.Lock(output),
			zap.LevelEnablerFunc(levelEnabler),
		)
		if c.RegisterHooks != nil {
			for _, h := range c.RegisterHooks {
				newCore = zapcore.RegisterHooks(newCore, h)
			}
		}
		cores = append(cores, newCore)
	}

	Logger = zap.New(
		zapcore.NewTee(cores...),
		zap.AddCaller(),
	)
}

func GetDefaultEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func Close() {
	Logger.Sync()
	if logFile != nil {
		logFile.Close()
	}
}
