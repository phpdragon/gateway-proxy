package eureka_client

import (
	"go.uber.org/zap"
	"log"
	"os"
)

type ClientLogger struct {
	zapLogger *zap.Logger
	baseLog   *log.Logger
}

// New creates a new ClientLogger Agent
func NewLogAgent(zapLogger *zap.Logger) *ClientLogger {
	return &ClientLogger{zapLogger: zapLogger, baseLog: log.New(os.Stderr, "", log.LstdFlags)}
}

func (log *ClientLogger) Debug(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Debug(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *ClientLogger) Info(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Info(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *ClientLogger) Warn(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Warn(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *ClientLogger) Error(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Error(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *ClientLogger) Panic(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Panic(msg)
	} else {
		log.baseLog.Panic(msg)
	}
}

func (log *ClientLogger) DPanic(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.DPanic(msg)
	} else {
		log.baseLog.Println(msg)
	}
}

func (log *ClientLogger) Fatal(msg string) {
	if nil != log.zapLogger {
		log.zapLogger.Fatal(msg)
	} else {
		log.baseLog.Fatal(msg)
	}
}
