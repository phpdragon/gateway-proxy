package config

import (
	"github.com/phpdragon/gateway-proxy/internal/components/logger"
)

func init() {
	logConfig := GetLogConfig()
	filename := logConfig.GetLogFilePath()
	logger.InitLog(filename)
}
