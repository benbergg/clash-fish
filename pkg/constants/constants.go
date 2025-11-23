package constants

import (
	"os"
	"path/filepath"
)

const (
	// AppName 应用名称
	AppName = "clash-fish"

	// Version 版本号
	Version = "0.1.0"

	// DefaultConfigFileName 默认配置文件名
	DefaultConfigFileName = "config.yaml"

	// DefaultPIDFileName PID 文件名
	DefaultPIDFileName = "clash-fish.pid"

	// DefaultLogFileName 日志文件名
	DefaultLogFileName = "clash-fish.log"
)

// GetDefaultConfigDir 获取默认配置目录
func GetDefaultConfigDir() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "/tmp/clash-fish"
	}
	return filepath.Join(homeDir, ".config", AppName)
}

// GetDefaultLogDir 获取默认日志目录
func GetDefaultLogDir() string {
	return filepath.Join(GetDefaultConfigDir(), "logs")
}

// GetDefaultProfilesDir 获取默认 profiles 目录
func GetDefaultProfilesDir() string {
	return filepath.Join(GetDefaultConfigDir(), "profiles")
}
