#!/bin/bash
# 检查 Clash-Fish 运行状态的脚本
# 在服务启动后，在另一个终端运行此脚本

GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo "========================================="
echo "  Clash-Fish 运行状态检查"
echo "========================================="
echo ""

# 1. 检查服务状态
echo -e "${BLUE}1. 检查服务状态${NC}"
echo "-----------------------------------"
./build/clash-fish status
echo ""

# 2. 检查 PID 文件
echo -e "${BLUE}2. 检查 PID 文件${NC}"
echo "-----------------------------------"
if [ -f ~/.config/clash-fish/clash-fish.pid ]; then
    PID=$(cat ~/.config/clash-fish/clash-fish.pid)
    echo -e "${GREEN}✓ PID 文件存在: $PID${NC}"

    # 检查进程是否存在
    if ps -p $PID > /dev/null; then
        echo -e "${GREEN}✓ 进程运行中 (PID: $PID)${NC}"
    else
        echo -e "${RED}✗ 进程不存在 (PID: $PID)${NC}"
    fi
else
    echo -e "${RED}✗ PID 文件不存在${NC}"
fi
echo ""

# 3. 检查 TUN 设备
echo -e "${BLUE}3. 检查 TUN 设备${NC}"
echo "-----------------------------------"
if ifconfig | grep -q "utun5"; then
    echo -e "${GREEN}✓ utun5 设备已创建${NC}"
    ifconfig utun5 | grep "inet"
else
    echo -e "${YELLOW}⚠ utun5 设备未找到${NC}"
    echo "  可能 TUN 模式未启用或使用其他设备名"
fi
echo ""

# 4. 检查路由表
echo -e "${BLUE}4. 检查路由表${NC}"
echo "-----------------------------------"
echo "查找 198.18.0.1 相关路由:"
netstat -nr | grep "198.18" || echo -e "${YELLOW}⚠ 未找到代理路由${NC}"
echo ""

# 5. 检查端口监听
echo -e "${BLUE}5. 检查端口监听${NC}"
echo "-----------------------------------"
echo "HTTP 端口 (7890):"
if lsof -i :7890 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 端口 7890 已监听${NC}"
    lsof -i :7890 | head -2
else
    echo -e "${YELLOW}⚠ 端口 7890 未监听${NC}"
fi
echo ""

echo "SOCKS5 端口 (7891):"
if lsof -i :7891 > /dev/null 2>&1; then
    echo -e "${GREEN}✓ 端口 7891 已监听${NC}"
    lsof -i :7891 | head -2
else
    echo -e "${YELLOW}⚠ 端口 7891 未监听${NC}"
fi
echo ""

# 6. 测试代理连接
echo -e "${BLUE}6. 测试代理连接${NC}"
echo "-----------------------------------"
echo "测试 HTTP 代理 (www.google.com)..."
if curl -x http://127.0.0.1:7890 -I -s --max-time 5 https://www.google.com > /dev/null 2>&1; then
    echo -e "${GREEN}✓ HTTP 代理工作正常${NC}"
else
    echo -e "${YELLOW}⚠ HTTP 代理测试失败（可能代理服务器未配置或不可用）${NC}"
fi
echo ""

# 7. 检查 VPN 共存
echo -e "${BLUE}7. 检查 VPN 共存${NC}"
echo "-----------------------------------"
if ifconfig | grep -q "utun4"; then
    echo -e "${GREEN}✓ VPN 连接存在 (utun4)${NC}"
    echo "VPN 路由:"
    netstat -nr | grep "utun4" | head -5
    echo ""
    echo -e "${BLUE}测试 VPN 路由优先级...${NC}"
    echo "10.8.0.0/24 路由:"
    netstat -nr | grep "10.8" || echo "未找到"
else
    echo -e "${YELLOW}⚠ VPN 未连接${NC}"
fi
echo ""

# 8. 检查日志
echo -e "${BLUE}8. 检查日志${NC}"
echo "-----------------------------------"
if [ -f ~/.config/clash-fish/logs/clash-fish.log ]; then
    echo "最近的日志:"
    tail -10 ~/.config/clash-fish/logs/clash-fish.log
else
    echo -e "${YELLOW}⚠ 日志文件不存在${NC}"
fi
echo ""

echo "========================================="
echo "  检查完成"
echo "========================================="
