#!/bin/bash
# Clash-Fish 手动测试脚本
# 需要在终端中运行并输入 sudo 密码

set -e

BINARY="./build/clash-fish"
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "  Clash-Fish 手动测试脚本"
echo "========================================="
echo ""

# 检查二进制文件
if [ ! -f "$BINARY" ]; then
    echo -e "${RED}✗ 错误: 找不到 $BINARY${NC}"
    echo "请先运行: make build"
    exit 1
fi

echo -e "${GREEN}✓ 找到二进制文件${NC}"
echo ""

# 测试 1: 配置验证
echo "========================================="
echo "测试 1: 配置验证"
echo "========================================="
$BINARY config validate
echo -e "${GREEN}✓ 配置验证通过${NC}"
echo ""

# 测试 2: 状态检查（启动前）
echo "========================================="
echo "测试 2: 状态检查（启动前）"
echo "========================================="
$BINARY status
echo ""

# 测试 3: 启动服务
echo "========================================="
echo "测试 3: 启动服务（需要 sudo 密码）"
echo "========================================="
echo -e "${YELLOW}提示: 服务将在前台运行，按 Ctrl+C 停止${NC}"
echo ""
read -p "准备好启动服务了吗? (y/n) " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]; then
    echo "测试已取消"
    exit 0
fi

echo ""
echo "正在启动服务..."
echo "================================"
sudo $BINARY start

# 注意：如果执行到这里，说明服务已经被 Ctrl+C 停止

echo ""
echo "========================================="
echo "测试 4: 验证服务已停止"
echo "========================================="
$BINARY status
echo ""

echo -e "${GREEN}✓ 所有测试完成${NC}"
