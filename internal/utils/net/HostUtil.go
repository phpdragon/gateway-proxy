package net

import (
	"net"
	"regexp"
)

// IsIpAddressPort 是否ip:port格式
func IsIpAddressPort(ipStr string) bool {
	if isOk, _ := regexp.MatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)(:\d+)$`, ipStr); isOk {
		return isOk
	}
	return false
}

// GetLocalIp 获取本地ip
func GetLocalIp() string {
	addr, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, address := range addr {
		// check the address type and if it is not a loopback the display it
		if ipNet, ok := address.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}
