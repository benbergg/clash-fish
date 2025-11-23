package main

import (
	"fmt"

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

	// TODO: 实现启动逻辑
	// 1. 检查服务是否已运行
	// 2. 加载配置文件
	// 3. 启动 mihomo 引擎
	// 4. 保存 PID 文件

	fmt.Println("✓ Clash-Fish started successfully")
	logger.Info().Msg("Service started")

	return nil
}

func init() {
	rootCmd.AddCommand(startCmd)
}
