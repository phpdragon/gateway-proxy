package utils

import (
	"os"
	"path/filepath"
	"strings"
)

/*
获取程序运行路径
*/
func GetCurrentDir() string {
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	return strings.Replace(dir, "\\", "/", -1)
}
