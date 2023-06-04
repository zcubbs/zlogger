# zlogger

This is an interface for Loggers in Go. It is inspired by [go-logger](https://github.com/arun0009/go-logger). 

## Usage

### Install

```bash
go get github.com/zcubbs/zlogger
```

### Logrus

[logrus](https://github.com/sirupsen/logrus)

```go
package main

import (
	"os"

	"github.com/zcubbs/zlogger/pkg/logger"
	"github.com/sirupsen/logrus"
)

func main() {
	logrusLog := logrus.New()
	logrusLog.SetFormatter(&logrus.JSONFormatter{})
	logrusLog.SetOutput(os.Stdout)
	logrusLog.SetLevel(logrus.DebugLevel)
	log, _ := logger.NewLogrusLogger(logrusLog)
	logger.ReplaceGlobals(log)
        //anywhere in your code you can now use logger.L() as its globally set
	logger.L().WithFields(logger.Fields{
		"foo": "bar",
	}).Info("direct")
}
```

### Zap

[zap](https://github.com/uber-go/zap)

```go
package main

import (
	"os"

	"github.com/zcubbs/zlogger/pkg/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	consoleEncoder := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	core := zapcore.NewCore(consoleEncoder,
		zapcore.Lock(zapcore.AddSync(os.Stderr)),
		zapcore.DebugLevel)
	zapLogger := zap.New(core)
	log, _ := logger.NewZapLogger(zapLogger)
 	logger.ReplaceGlobals(log)
        //anywhere in your code you can now use logger.L() as its globally set
	logger.L().WithFields(logger.Fields{
		"foo": "bar",
	}).Info("direct")
}
```

## License

MIT
