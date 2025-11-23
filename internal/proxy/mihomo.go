package proxy

import (
	"fmt"
	"os"

	"github.com/metacubex/mihomo/config"
	"github.com/metacubex/mihomo/hub/executor"
	"github.com/metacubex/mihomo/log"
)

// MihomoEngine Mihomo 引擎封装
type MihomoEngine struct {
	configPath string
	homeDir    string
	running    bool
}

// NewMihomoEngine 创建 Mihomo 引擎实例
func NewMihomoEngine(configPath, homeDir string) *MihomoEngine {
	return &MihomoEngine{
		configPath: configPath,
		homeDir:    homeDir,
		running:    false,
	}
}

// Start 启动 Mihomo 引擎
func (e *MihomoEngine) Start() error {
	if e.running {
		return fmt.Errorf("mihomo engine is already running")
	}

	// 设置 mihomo 的日志级别
	log.SetLevel(log.INFO)

	// 读取配置文件
	configData, err := os.ReadFile(e.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	cfg, err := config.Parse(configData)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 应用配置
	executor.ApplyConfig(cfg, true)

	e.running = true
	log.Infoln("Mihomo engine started successfully")

	return nil
}

// Stop 停止 Mihomo 引擎
func (e *MihomoEngine) Stop() error {
	if !e.running {
		return fmt.Errorf("mihomo engine is not running")
	}

	// mihomo 的优雅关闭
	// 注意：mihomo 库没有提供直接的 Stop 方法
	// 实际的停止会通过进程管理来实现

	e.running = false
	log.Infoln("Mihomo engine stopped")

	return nil
}

// IsRunning 检查引擎是否运行中
func (e *MihomoEngine) IsRunning() bool {
	return e.running
}

// Reload 重新加载配置
func (e *MihomoEngine) Reload() error {
	if !e.running {
		return fmt.Errorf("mihomo engine is not running")
	}

	// 读取配置文件
	configData, err := os.ReadFile(e.configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	// 解析配置
	cfg, err := config.Parse(configData)
	if err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	// 重新应用配置
	executor.ApplyConfig(cfg, false)

	log.Infoln("Mihomo engine reloaded")

	return nil
}
