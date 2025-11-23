package main

import (
	"fmt"

	"github.com/clash-fish/clash-fish/internal/proxy"
	"github.com/clash-fish/clash-fish/pkg/logger"
	"github.com/clash-fish/clash-fish/pkg/utils"
	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop clash-fish service",
	Long:  `Stop the running clash-fish service and cleanup system settings.`,
	RunE:  runStop,
}

func runStop(cmd *cobra.Command, args []string) error {
	// 检查 root 权限
	if err := utils.CheckRoot(); err != nil {
		return err
	}

	logger.Info().Msg("Stopping clash-fish service...")

	// 创建代理管理器
	manager := proxy.NewManager(configDir)

	// 停止服务
	if err := manager.Stop(); err != nil {
		return fmt.Errorf("failed to stop service: %w", err)
	}

	fmt.Println("✓ Clash-Fish stopped successfully")
	logger.Info().Msg("Service stopped")

	return nil
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
