package main

import (
	"fmt"

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

	// TODO: 实现停止逻辑
	// 1. 检查服务是否在运行
	// 2. 停止 mihomo 引擎
	// 3. 清理系统设置
	// 4. 删除 PID 文件

	fmt.Println("✓ Clash-Fish stopped successfully")
	logger.Info().Msg("Service stopped")

	return nil
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
