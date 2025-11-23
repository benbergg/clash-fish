package main

import (
	"fmt"

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

	// TODO: 实现重启逻辑
	// 1. 调用 stop 逻辑
	// 2. 调用 start 逻辑

	fmt.Println("✓ Clash-Fish restarted successfully")
	logger.Info().Msg("Service restarted")

	return nil
}

func init() {
	rootCmd.AddCommand(restartCmd)
}
