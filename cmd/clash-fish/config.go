package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/clash-fish/clash-fish/internal/config"
	"github.com/clash-fish/clash-fish/pkg/logger"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration files",
	Long:  `Manage clash-fish configuration files including init, edit, validate and show.`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration",
	Long:  `Create default configuration file in the config directory.`,
	RunE:  runConfigInit,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration file",
	Long:  `Open the configuration file in default editor.`,
	RunE:  runConfigEdit,
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate configuration file",
	Long:  `Check if the configuration file is valid.`,
	RunE:  runConfigValidate,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current configuration file content.`,
	RunE:  runConfigShow,
}

func runConfigInit(cmd *cobra.Command, args []string) error {
	logger.Info().Str("dir", configDir).Msg("Initializing configuration...")

	// 创建配置管理器
	mgr := config.NewManager(configDir)

	// 检查配置是否已存在
	if mgr.Exists() {
		fmt.Printf("⚠ Configuration already exists at: %s\n", mgr.GetConfigPath())
		fmt.Println("Use 'clash-fish config show' to view current configuration")
		return nil
	}

	// 初始化配置
	if err := mgr.Init(); err != nil {
		return fmt.Errorf("failed to initialize configuration: %w", err)
	}

	fmt.Printf("✓ Configuration initialized at: %s\n", configDir)
	fmt.Printf("  Config file: %s\n", mgr.GetConfigPath())
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Edit config: clash-fish config edit")
	fmt.Println("  2. Or add subscription: clash-fish profile add <name> <url>")
	fmt.Println("  3. Start service: sudo clash-fish start")

	logger.Info().Msg("Configuration initialized successfully")

	return nil
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	mgr := config.NewManager(configDir)

	// 检查配置文件是否存在
	if !mgr.Exists() {
		return fmt.Errorf("configuration not found, run 'clash-fish config init' first")
	}

	// 获取编辑器
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
	}

	// 打开编辑器
	configPath := mgr.GetConfigPath()
	fmt.Printf("Opening %s with %s...\n", configPath, editor)

	editorCmd := exec.Command(editor, configPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	fmt.Println("✓ Configuration edited")

	return nil
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	logger.Info().Msg("Validating configuration...")

	mgr := config.NewManager(configDir)

	// 检查配置文件是否存在
	if !mgr.Exists() {
		return fmt.Errorf("configuration not found, run 'clash-fish config init' first")
	}

	// 加载配置
	cfg, err := mgr.Load()
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// 验证配置
	if err := mgr.Validate(cfg); err != nil {
		fmt.Println("✗ Configuration is invalid:")
		fmt.Printf("  Error: %v\n", err)
		return err
	}

	fmt.Println("✓ Configuration is valid")
	fmt.Printf("  Mode: %s\n", cfg.Mode)
	fmt.Printf("  HTTP Port: %d\n", cfg.Port)
	fmt.Printf("  SOCKS Port: %d\n", cfg.SocksPort)
	fmt.Printf("  TUN Enabled: %v\n", cfg.TUN.Enable)
	fmt.Printf("  DNS Enabled: %v\n", cfg.DNS.Enable)

	logger.Info().Msg("Configuration validated successfully")

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	mgr := config.NewManager(configDir)

	// 检查配置文件是否存在
	if !mgr.Exists() {
		return fmt.Errorf("configuration not found, run 'clash-fish config init' first")
	}

	// 读取配置文件
	data, err := os.ReadFile(mgr.GetConfigPath())
	if err != nil {
		return fmt.Errorf("failed to read configuration: %w", err)
	}

	// 解析配置以便格式化显示
	var cfg config.Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	fmt.Println("=== Current Configuration ===")
	fmt.Printf("Config File: %s\n", mgr.GetConfigPath())
	fmt.Printf("Config Dir:  %s\n\n", configDir)

	// 显示主要配置
	fmt.Printf("Mode:        %s\n", cfg.Mode)
	fmt.Printf("HTTP Port:   %d\n", cfg.Port)
	fmt.Printf("SOCKS Port:  %d\n", cfg.SocksPort)
	fmt.Printf("Log Level:   %s\n", cfg.LogLevel)
	fmt.Printf("Allow LAN:   %v\n\n", cfg.AllowLan)

	// TUN 配置
	fmt.Printf("TUN Mode:    %v\n", cfg.TUN.Enable)
	if cfg.TUN.Enable {
		fmt.Printf("  Stack:     %s\n", cfg.TUN.Stack)
		fmt.Printf("  Auto Route: %v\n", cfg.TUN.AutoRoute)
	}

	// DNS 配置
	fmt.Printf("\nDNS:         %v\n", cfg.DNS.Enable)
	if cfg.DNS.Enable {
		fmt.Printf("  Mode:      %s\n", cfg.DNS.EnhancedMode)
		fmt.Printf("  Listen:    %s\n", cfg.DNS.Listen)
	}

	// 代理配置
	fmt.Printf("\nProxies:     %d configured\n", len(cfg.Proxies))
	fmt.Printf("Proxy Groups: %d configured\n", len(cfg.ProxyGroups))
	fmt.Printf("Rules:       %d configured\n", len(cfg.Rules))

	fmt.Println("\nUse 'clash-fish config edit' to modify configuration")

	return nil
}

func init() {
	// 添加子命令
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configShowCmd)

	// 添加到根命令
	rootCmd.AddCommand(configCmd)
}
