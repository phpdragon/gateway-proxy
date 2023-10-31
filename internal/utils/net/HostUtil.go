package net

import "regexp"

// IsIpAddressPort 是否ip:port格式
func IsIpAddressPort(ipStr string) bool {
	if isOk, _ := regexp.MatchString(`^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)(:\d+)$`, ipStr); isOk {
		return isOk
	}
	return false
}
