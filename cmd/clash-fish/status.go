package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show clash-fish service status",
	Long:  `Display the current status of clash-fish service including VPN detection.`,
	RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
	fmt.Println("=== Clash-Fish Status ===")

	// TODO: 实现状态查看逻辑
	// 1. 检查服务是否在运行
	// 2. 显示 VPN 连接状态
	// 3. 显示配置信息（模式、端口等）
	// 4. 显示流量统计（可选）

	fmt.Println("Service:    ✗ Not Running")
	fmt.Println("VPN:        - Not Detected")
	fmt.Println("Mode:       - N/A")
	fmt.Println("HTTP Port:  - N/A")
	fmt.Println("SOCKS Port: - N/A")

	return nil
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
