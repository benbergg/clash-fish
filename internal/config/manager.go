package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Manager 配置管理器
type Manager struct {
	configDir  string
	configPath string
	config     *Config
}

// NewManager 创建配置管理器
func NewManager(configDir string) *Manager {
	return &Manager{
		configDir:  configDir,
		configPath: filepath.Join(configDir, "config.yaml"),
	}
}

// Init 初始化配置目录和文件
func (m *Manager) Init() error {
	// 创建配置目录
	if err := os.MkdirAll(m.configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// 创建子目录
	dirs := []string{
		filepath.Join(m.configDir, "logs"),
		filepath.Join(m.configDir, "profiles"),
		filepath.Join(m.configDir, "cache"),
	}

	for _, dir := range dirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	// 如果配置文件不存在，创建默认配置
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		if err := m.CreateDefault(); err != nil {
			return fmt.Errorf("failed to create default config: %w", err)
		}
	}

	return nil
}

// CreateDefault 创建默认配置文件
func (m *Manager) CreateDefault() error {
	config := GetDefaultConfig()
	return m.Save(config)
}

// Load 加载配置文件
func (m *Manager) Load() (*Config, error) {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	m.config = &config
	return &config, nil
}

// Save 保存配置文件
func (m *Manager) Save(config *Config) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// 添加配置文件头部注释
	header := []byte("# Clash-Fish Configuration\n# Auto-generated configuration file\n\n")
	data = append(header, data...)

	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	m.config = config
	return nil
}

// Validate 验证配置文件
func (m *Manager) Validate(config *Config) error {
	// 验证必需字段
	if config.Port <= 0 || config.Port > 65535 {
		return fmt.Errorf("invalid port: %d", config.Port)
	}

	if config.SocksPort <= 0 || config.SocksPort > 65535 {
		return fmt.Errorf("invalid socks-port: %d", config.SocksPort)
	}

	// 验证模式
	validModes := map[string]bool{
		"rule":   true,
		"global": true,
		"direct": true,
	}
	if !validModes[config.Mode] {
		return fmt.Errorf("invalid mode: %s (must be rule/global/direct)", config.Mode)
	}

	// 验证日志级别
	validLogLevels := map[string]bool{
		"info":    true,
		"warning": true,
		"error":   true,
		"debug":   true,
		"silent":  true,
	}
	if !validLogLevels[config.LogLevel] {
		return fmt.Errorf("invalid log-level: %s", config.LogLevel)
	}

	// 验证 TUN 配置
	if config.TUN.Enable {
		validStacks := map[string]bool{
			"system": true,
			"gvisor": true,
		}
		if !validStacks[config.TUN.Stack] {
			return fmt.Errorf("invalid tun.stack: %s (must be system/gvisor)", config.TUN.Stack)
		}
	}

	// 验证 DNS 配置
	if config.DNS.Enable {
		validModes := map[string]bool{
			"fake-ip":    true,
			"redir-host": true,
		}
		if !validModes[config.DNS.EnhancedMode] {
			return fmt.Errorf("invalid dns.enhanced-mode: %s", config.DNS.EnhancedMode)
		}
	}

	return nil
}

// GetConfigPath 获取配置文件路径
func (m *Manager) GetConfigPath() string {
	return m.configPath
}

// GetConfigDir 获取配置目录路径
func (m *Manager) GetConfigDir() string {
	return m.configDir
}

// Exists 检查配置文件是否存在
func (m *Manager) Exists() bool {
	_, err := os.Stat(m.configPath)
	return err == nil
}
