package main

import (
	"fmt"

	"github.com/clash-fish/clash-fish/pkg/constants"
	"github.com/spf13/cobra"
)

var (
	logLines int
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "View logs",
	Long:  `View clash-fish service logs.`,
}

var logTailCmd = &cobra.Command{
	Use:   "tail",
	Short: "Show real-time logs",
	Long:  `Display real-time logs from clash-fish service (like tail -f).`,
	RunE:  runLogTail,
}

var logShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show historical logs",
	Long:  `Display historical logs from clash-fish service.`,
	RunE:  runLogShow,
}

func runLogTail(cmd *cobra.Command, args []string) error {
	logDir := constants.GetDefaultLogDir()
	logFile := fmt.Sprintf("%s/%s", logDir, constants.DefaultLogFileName)

	fmt.Printf("Tailing logs from: %s\n", logFile)
	fmt.Println("Press Ctrl+C to stop...")
	fmt.Println("---")

	// TODO: 实现实时日志查看逻辑
	// 1. 打开日志文件
	// 2. 使用 tail -f 或类似机制
	// 3. 实时显示新日志

	fmt.Println("(Real-time logs will be displayed here)")

	return nil
}

func runLogShow(cmd *cobra.Command, args []string) error {
	logDir := constants.GetDefaultLogDir()
	logFile := fmt.Sprintf("%s/%s", logDir, constants.DefaultLogFileName)

	fmt.Printf("Showing logs from: %s\n", logFile)
	fmt.Println("---")

	// TODO: 实现历史日志查看逻辑
	// 1. 读取日志文件
	// 2. 显示最后 N 行（默认 50 行）

	fmt.Println("(Historical logs will be displayed here)")

	return nil
}

func init() {
	// 添加标志
	logShowCmd.Flags().IntVarP(&logLines, "lines", "n", 50, "number of lines to show")

	// 添加子命令
	logCmd.AddCommand(logTailCmd)
	logCmd.AddCommand(logShowCmd)

	// 添加到根命令
	rootCmd.AddCommand(logCmd)
}
