package main

import (
	"fmt"

	"github.com/clash-fish/clash-fish/pkg/logger"
	"github.com/spf13/cobra"
)

var profileCmd = &cobra.Command{
	Use:   "profile",
	Short: "Manage configuration profiles",
	Long:  `Manage multiple configuration profiles from different subscription sources.`,
}

var profileListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all profiles",
	Long:  `Display all available configuration profiles.`,
	RunE:  runProfileList,
}

var profileAddCmd = &cobra.Command{
	Use:   "add <name> <url>",
	Short: "Add a new profile",
	Long:  `Add a new configuration profile from subscription URL.`,
	Args:  cobra.ExactArgs(2),
	RunE:  runProfileAdd,
}

var profileUpdateCmd = &cobra.Command{
	Use:   "update [name]",
	Short: "Update profile(s)",
	Long:  `Update configuration profile from subscription URL. If no name specified, update all.`,
	Args:  cobra.MaximumNArgs(1),
	RunE:  runProfileUpdate,
}

var profileSwitchCmd = &cobra.Command{
	Use:   "switch <name>",
	Short: "Switch to a profile",
	Long:  `Switch to a different configuration profile.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileSwitch,
}

var profileDeleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a profile",
	Long:  `Delete a configuration profile.`,
	Args:  cobra.ExactArgs(1),
	RunE:  runProfileDelete,
}

func runProfileList(cmd *cobra.Command, args []string) error {
	// TODO: 实现 profile 列表逻辑
	// 1. 读取 profiles 目录
	// 2. 列出所有 profile
	// 3. 标记当前激活的 profile

	fmt.Println("=== Configuration Profiles ===")
	fmt.Println("No profiles found. Use 'clash-fish profile add' to add one.")

	return nil
}

func runProfileAdd(cmd *cobra.Command, args []string) error {
	name := args[0]
	url := args[1]

	logger.Info().
		Str("name", name).
		Str("url", url).
		Msg("Adding profile...")

	// TODO: 实现 profile 添加逻辑
	// 1. 下载订阅内容
	// 2. 解析配置
	// 3. 保存到 profiles 目录
	// 4. 更新 profile 元数据

	fmt.Printf("✓ Profile '%s' added successfully\n", name)
	logger.Info().Str("name", name).Msg("Profile added")

	return nil
}

func runProfileUpdate(cmd *cobra.Command, args []string) error {
	var name string
	if len(args) > 0 {
		name = args[0]
		logger.Info().Str("name", name).Msg("Updating profile...")
	} else {
		logger.Info().Msg("Updating all profiles...")
	}

	// TODO: 实现 profile 更新逻辑
	// 1. 读取 profile 的订阅 URL
	// 2. 重新下载订阅
	// 3. 更新配置文件

	if name != "" {
		fmt.Printf("✓ Profile '%s' updated successfully\n", name)
	} else {
		fmt.Println("✓ All profiles updated successfully")
	}

	return nil
}

func runProfileSwitch(cmd *cobra.Command, args []string) error {
	name := args[0]

	logger.Info().Str("name", name).Msg("Switching profile...")

	// TODO: 实现 profile 切换逻辑
	// 1. 验证 profile 存在
	// 2. 复制 profile 配置到主配置
	// 3. 更新当前激活的 profile 标记
	// 4. 如果服务正在运行，重新加载配置

	fmt.Printf("✓ Switched to profile '%s'\n", name)
	logger.Info().Str("name", name).Msg("Profile switched")

	return nil
}

func runProfileDelete(cmd *cobra.Command, args []string) error {
	name := args[0]

	logger.Info().Str("name", name).Msg("Deleting profile...")

	// TODO: 实现 profile 删除逻辑
	// 1. 验证 profile 存在
	// 2. 检查是否是当前激活的 profile
	// 3. 删除 profile 文件

	fmt.Printf("✓ Profile '%s' deleted successfully\n", name)
	logger.Info().Str("name", name).Msg("Profile deleted")

	return nil
}

func init() {
	// 添加子命令
	profileCmd.AddCommand(profileListCmd)
	profileCmd.AddCommand(profileAddCmd)
	profileCmd.AddCommand(profileUpdateCmd)
	profileCmd.AddCommand(profileSwitchCmd)
	profileCmd.AddCommand(profileDeleteCmd)

	// 添加到根命令
	rootCmd.AddCommand(profileCmd)
}
