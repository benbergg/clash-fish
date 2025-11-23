package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/clash-fish/clash-fish/internal/proxy"
	"github.com/clash-fish/clash-fish/pkg/logger"
	"github.com/clash-fish/clash-fish/pkg/utils"
	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start clash-fish service",
	Long:  `Start the clash-fish transparent proxy service with TUN mode.`,
	RunE:  runStart,
}

func runStart(cmd *cobra.Command, args []string) error {
	// 检查 root 权限
	if err := utils.CheckRoot(); err != nil {
		return err
	}

	logger.Info().Msg("Starting clash-fish service...")

	// 创建代理管理器
	manager := proxy.NewManager(configDir)

	// 启动服务
	if err := manager.Start(); err != nil {
		return fmt.Errorf("failed to start service: %w", err)
	}

	fmt.Println("✓ Clash-Fish started successfully")
	fmt.Printf("  Config: %s\n", manager.GetConfigPath())
	fmt.Println("\nService is running in foreground. Press Ctrl+C to stop.")

	// 设置信号处理
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// 等待退出信号
	sig := <-sigCh
	logger.Info().Str("signal", sig.String()).Msg("Received signal, shutting down...")

	fmt.Println("\nStopping Clash-Fish...")

	// 停止服务
	if err := manager.Stop(); err != nil {
		logger.Error().Err(err).Msg("Failed to stop service gracefully")
		return err
	}

	fmt.Println("✓ Clash-Fish stopped")

	return nil
}

func init() {
	rootCmd.AddCommand(startCmd)
}
