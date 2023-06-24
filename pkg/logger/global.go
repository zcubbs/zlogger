package logger

import (
	"github.com/sirupsen/logrus"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"sync"
)

var (
	_globalMU sync.RWMutex
	_globalL  Logger
)

// ReplaceGlobals replaces the global Logger and returns a
// function to restore the original values. It's safe for concurrent use.
func ReplaceGlobals(logger Logger) func() {
	_globalMU.Lock()
	prev := _globalL
	_globalL = logger
	_globalMU.Unlock()
	return func() { ReplaceGlobals(prev) }
}

// L returns the global Logger, which can be reconfigured with ReplaceGlobals.
// It's safe for concurrent use.
func L() Logger {
	_globalMU.RLock()
	l := _globalL
	_globalMU.RUnlock()
	return l
}

type Type string

const (
	LogrusLogger Type = "logrus"
	ZapLogger    Type = "zap"
)

func SetupLogger(loggerType Type) {
	if loggerType == LogrusLogger {
		setupLogrus()
		return
	}
	if loggerType == ZapLogger {
		setupZap()
		return
	}

	panic("invalid logger type")
}

func setupLogrus() {
	logrusLog := logrus.New()
	logrusLog.SetReportCaller(false)
	logrusLog.SetFormatter(&logrus.JSONFormatter{})
	logrusLog.SetOutput(os.Stdout)
	logrusLog.SetLevel(logrus.InfoLevel)
	log, _ := NewLogrusLogger(logrusLog)
	ReplaceGlobals(log)
}

func setupZap() {
	consoleEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.EpochTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})
	core := zapcore.NewCore(consoleEncoder,
		zapcore.Lock(zapcore.AddSync(os.Stderr)),
		zapcore.DebugLevel)
	zapLogger := zap.New(core)
	log, _ := NewZapLogger(zapLogger)
	ReplaceGlobals(log)
}
