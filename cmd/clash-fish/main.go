package main

import (
	"fmt"
	"os"

	"github.com/clash-fish/clash-fish/pkg/constants"
	"github.com/clash-fish/clash-fish/pkg/logger"
	"github.com/spf13/cobra"
)

var (
	// 全局标志
	configDir string
	debug     bool
)

var rootCmd = &cobra.Command{
	Use:     constants.AppName,
	Short:   "A transparent proxy tool based on mihomo",
	Long:    `Clash-Fish is a CLI tool for macOS that provides transparent proxy with VPN coexistence.`,
	Version: constants.Version,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 初始化日志系统
		logDir := constants.GetDefaultLogDir()
		if err := logger.Init(logDir, debug); err != nil {
			return fmt.Errorf("failed to initialize logger: %w", err)
		}
		return nil
	},
}

func init() {
	// 全局标志
	rootCmd.PersistentFlags().StringVar(&configDir, "config-dir", constants.GetDefaultConfigDir(), "config directory")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug mode")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
