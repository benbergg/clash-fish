package utils

import (
	"fmt"
	"os"
)

// CheckRoot 检查是否有 root 权限
func CheckRoot() error {
	if !IsRoot() {
		return fmt.Errorf("this command requires root privileges, please run with sudo")
	}
	return nil
}

// IsRoot 判断是否是 root 用户
func IsRoot() bool {
	return os.Geteuid() == 0
}
