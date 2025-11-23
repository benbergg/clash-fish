package main

import (
	"fmt"

	"github.com/clash-fish/clash-fish/internal/proxy"
	"github.com/clash-fish/clash-fish/pkg/logger"
	"github.com/clash-fish/clash-fish/pkg/utils"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart clash-fish service",
	Long:  `Restart the clash-fish service (stop then start).`,
	RunE:  runRestart,
}

func runRestart(cmd *cobra.Command, args []string) error {
	// 检查 root 权限
	if err := utils.CheckRoot(); err != nil {
		return err
	}

	logger.Info().Msg("Restarting clash-fish service...")

	// 创建代理管理器
	manager := proxy.NewManager(configDir)

	// 重启服务
	if err := manager.Restart(); err != nil {
		return fmt.Errorf("failed to restart service: %w", err)
	}

	fmt.Println("✓ Clash-Fish restarted successfully")
	fmt.Printf("  Config: %s\n", manager.GetConfigPath())
	logger.Info().Msg("Service restarted")

	return nil
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
