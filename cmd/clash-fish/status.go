package main

import (
	"fmt"

	"github.com/clash-fish/clash-fish/internal/config"
	"github.com/clash-fish/clash-fish/internal/proxy"
	"github.com/clash-fish/clash-fish/internal/system"
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

	// 创建代理管理器
	manager := proxy.NewManager(configDir)

	// 检查服务状态
	if manager.IsRunning() {
		pid, _ := manager.GetPID()
		fmt.Printf("Service:    ✓ Running (PID: %d)\n", pid)
	} else {
		fmt.Println("Service:    ✗ Not Running")
	}

	// VPN 检测
	vpnInfo, err := system.DetectVPN()
	if err != nil {
		fmt.Printf("VPN:        ✗ Detection Failed: %v\n", err)
	} else if vpnInfo.Active {
		fmt.Printf("VPN:        ✓ Active (%s: %s)\n", vpnInfo.Interface, vpnInfo.IP)
	} else {
		fmt.Println("VPN:        - Not Detected")
	}

	// 配置信息
	cfgMgr := config.NewManager(configDir)
	if cfgMgr.Exists() {
		cfg, err := cfgMgr.Load()
		if err != nil {
			fmt.Printf("Config:     ✗ Load Failed: %v\n", err)
		} else {
			fmt.Printf("Mode:       %s\n", cfg.Mode)
			fmt.Printf("HTTP Port:  %d\n", cfg.Port)
			fmt.Printf("SOCKS Port: %d\n", cfg.SocksPort)
			fmt.Printf("TUN Mode:   %v\n", cfg.TUN.Enable)
			if cfg.TUN.Enable {
				fmt.Printf("  Stack:    %s\n", cfg.TUN.Stack)
			}
		}
	} else {
		fmt.Println("Config:     ✗ Not Initialized")
		fmt.Println("            Run 'clash-fish config init' to create configuration")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
