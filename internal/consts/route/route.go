package route

const SysRefreshKey = "6b71e64a1081cabf8fdc228bc0cbc74d"

const RouterSystem = "/system/"

// RspModeDefault 明文应答
const RspModeDefault = 0

// RspModeEncrypt 加密应答
const RspModeEncrypt = 1

// CrossModeDefault 不处理跨域(由下游处理)
const CrossModeDefault = 0

// CrossModeAllow 允许跨域
const CrossModeAllow = 1

// CrossModeConfig 配置跨域
const CrossModeConfig = 2
