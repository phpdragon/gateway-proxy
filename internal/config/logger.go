package config

import (
	"fmt"
	"github.com/phpdragon/gateway-proxy/internal/utils/date"
	"go.uber.org/zap"
	zapCore "go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
)

var logger *zap.SugaredLogger

func NewLogger() {
	var filePath = getLogFilePath(appConfig.Log.Path)
	initLog(filePath)
}

func initLog(filename string) {
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

	// High-priority output should also go to standard errorcode, and low-priority
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

	log := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	logger = log.Sugar()

	defer logger.Sync()
}

func Logger() *zap.SugaredLogger {
	return logger
}

func GetAtomicLevel() zap.AtomicLevel {
	return zap.NewAtomicLevel()
}

// GetLogFilePath 获取日志文件路径
func getLogFilePath(filename string) string {
	path := strings.TrimRight(filename, "/")
	return fmt.Sprintf("%s/%s_%s.log", path, appConfig.AppName, date.GetDatetimeYmd())
}
