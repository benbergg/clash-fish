package proxy

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"syscall"

	"github.com/clash-fish/clash-fish/internal/system"
	"github.com/clash-fish/clash-fish/pkg/logger"
)

// Manager 代理管理器
type Manager struct {
	engine     *MihomoEngine
	configPath string
	homeDir    string
	pidFile    string
}

// NewManager 创建代理管理器
func NewManager(homeDir string) *Manager {
	return &Manager{
		configPath: filepath.Join(homeDir, "config.yaml"),
		homeDir:    homeDir,
		pidFile:    filepath.Join(homeDir, "clash-fish.pid"),
	}
}

// Start 启动服务
func (m *Manager) Start() error {
	// 检查是否已运行
	if m.IsRunning() {
		return fmt.Errorf("service is already running (PID file exists)")
	}

	// 检查配置文件是否存在
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		return fmt.Errorf("configuration file not found: %s\nPlease run 'clash-fish config init' first", m.configPath)
	}

	// VPN 检测
	vpnInfo, err := system.DetectVPN()
	if err != nil {
		logger.Warn().Err(err).Msg("Failed to detect VPN")
	} else if vpnInfo.Active {
		logger.Info().
			Str("interface", vpnInfo.Interface).
			Str("ip", vpnInfo.IP).
			Str("network", vpnInfo.Network).
			Msg("VPN detected, proxy will coexist with VPN")
	}

	// 创建 Mihomo 引擎
	m.engine = NewMihomoEngine(m.configPath, m.homeDir)

	// 启动引擎
	if err := m.engine.Start(); err != nil {
		return fmt.Errorf("failed to start mihomo engine: %w", err)
	}

	// 保存 PID 文件
	if err := m.savePID(); err != nil {
		// 启动失败，停止引擎
		m.engine.Stop()
		return fmt.Errorf("failed to save PID file: %w", err)
	}

	logger.Info().
		Str("config", m.configPath).
		Str("pid_file", m.pidFile).
		Msg("Service started successfully")

	return nil
}

// Stop 停止服务
func (m *Manager) Stop() error {
	// 检查是否在运行
	if !m.IsRunning() {
		return fmt.Errorf("service is not running")
	}

	// 读取 PID
	pid, err := m.readPID()
	if err != nil {
		return fmt.Errorf("failed to read PID: %w", err)
	}

	// 检查进程是否存在
	process, err := os.FindProcess(pid)
	if err != nil {
		// PID 文件存在但进程不存在，清理 PID 文件
		os.Remove(m.pidFile)
		return fmt.Errorf("process not found, PID file cleaned")
	}

	// 发送 SIGTERM 信号
	if err := process.Signal(syscall.SIGTERM); err != nil {
		logger.Warn().Err(err).Msg("Failed to send SIGTERM, trying SIGKILL")
		// 如果 SIGTERM 失败，尝试 SIGKILL
		if err := process.Kill(); err != nil {
			return fmt.Errorf("failed to kill process: %w", err)
		}
	}

	// 删除 PID 文件
	if err := os.Remove(m.pidFile); err != nil {
		logger.Warn().Err(err).Msg("Failed to remove PID file")
	}

	logger.Info().Int("pid", pid).Msg("Service stopped successfully")

	return nil
}

// Restart 重启服务
func (m *Manager) Restart() error {
	// 如果正在运行，先停止
	if m.IsRunning() {
		if err := m.Stop(); err != nil {
			return fmt.Errorf("failed to stop service: %w", err)
		}
	}

	// 启动服务
	if err := m.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	return nil
}

// IsRunning 检查服务是否运行
func (m *Manager) IsRunning() bool {
	// 检查 PID 文件是否存在
	if _, err := os.Stat(m.pidFile); os.IsNotExist(err) {
		return false
	}

	// 读取 PID
	pid, err := m.readPID()
	if err != nil {
		return false
	}

	// 检查进程是否存在
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	// 尝试发送信号 0 检查进程是否真的存在
	err = process.Signal(syscall.Signal(0))
	return err == nil
}

// GetPID 获取运行中的 PID
func (m *Manager) GetPID() (int, error) {
	if !m.IsRunning() {
		return 0, fmt.Errorf("service is not running")
	}
	return m.readPID()
}

// savePID 保存 PID 到文件
func (m *Manager) savePID() error {
	pid := os.Getpid()
	content := []byte(strconv.Itoa(pid))
	return os.WriteFile(m.pidFile, content, 0644)
}

// readPID 从文件读取 PID
func (m *Manager) readPID() (int, error) {
	content, err := os.ReadFile(m.pidFile)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(string(content))
	if err != nil {
		return 0, fmt.Errorf("invalid PID in file: %w", err)
	}

	return pid, nil
}

// GetConfigPath 获取配置文件路径
func (m *Manager) GetConfigPath() string {
	return m.configPath
}
