package log

import (
	"go.uber.org/zap"
	zapCore "go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
)

// error logger
var errorLogger *zap.Logger

// InitLog TODO: 代码组织形式要优化一下，使用init进行初始化
func InitLog(filename string) {
	fileWriter := zapCore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    1024, // megabytes
		MaxBackups: 3,
		MaxAge:     7, //days
		LocalTime:  true,
		Compress:   true,
	})

	//encoder := zap.NewProductionEncoderConfig()
	//encoder.EncodeTime = zapCore.ISO8601TimeEncoder

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleDebugging := zapCore.Lock(os.Stdout)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	productEncoder := zapCore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapCore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapCore.Cores, then tee the four cores together.
	core := zapCore.NewTee(
		// 打印在控制台
		zapCore.NewCore(consoleEncoder, consoleDebugging, zap.NewAtomicLevelAt(zap.DebugLevel)),
		// 打印在文件中,仅打印Info级别以上的日志
		zapCore.NewCore(productEncoder, fileWriter, zap.NewAtomicLevelAt(zap.InfoLevel)),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	errorLogger = logger
	defer errorLogger.Sync()
}

func GetLogger() *zap.Logger {
	return errorLogger
}

func GetAtomicLevel() zap.AtomicLevel {
	return zap.NewAtomicLevel()
}

func Debug(msg string, fields ...zap.Field) {
	errorLogger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	errorLogger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	errorLogger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	errorLogger.Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	errorLogger.Panic(msg, fields...)
}

func DPanic(msg string, fields ...zap.Field) {
	errorLogger.DPanic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	errorLogger.Fatal(msg, fields...)
}
