# Clash-Fish 测试报告

## 测试环境
- **系统**: macOS 15.7.2 (Sequoia)
- **日期**: 2025-11-23
- **版本**: v0.1.0
- **VPN 状态**: Active (utun4: 10.8.0.105)

---

## 自动化测试结果

### ✅ 配置管理测试

#### 1. 配置初始化
```bash
$ ./build/clash-fish config init
✓ Configuration initialized at: /Users/lg/.config/clash-fish
  Config file: /Users/lg/.config/clash-fish/config.yaml
```
**结果**: PASS ✅

#### 2. 配置验证
```bash
$ ./build/clash-fish config validate
✓ Configuration is valid
  Mode: rule
  HTTP Port: 7890
  SOCKS Port: 7891
  TUN Enabled: true
  DNS Enabled: true
```
**结果**: PASS ✅

#### 3. 配置显示
```bash
$ ./build/clash-fish config show
=== Current Configuration ===
Config File: /Users/lg/.config/clash-fish/config.yaml
Mode:        rule
HTTP Port:   7890
SOCKS Port:  7891
TUN Mode:    true
  Stack:     system
  Auto Route: true
```
**结果**: PASS ✅

### ✅ VPN 检测测试

```bash
$ ./build/clash-fish status
VPN:        ✓ Active (utun4: 10.8.0.105)
```
**结果**: PASS ✅
- 成功检测到 OpenVPN 连接
- 正确识别 utun4 接口
- 正确解析私网 IP (10.8.0.105)

### ✅ 权限检查测试

```bash
$ ./build/clash-fish start
Error: this command requires root privileges, please run with sudo
```
**结果**: PASS ✅
- 正确检测非 root 用户
- 提示用户使用 sudo
- 防止未授权操作

---

## 需要手动测试的功能

由于以下功能需要 **sudo 权限和交互式终端**，需要用户手动测试：

### 🔧 手动测试步骤

#### 测试 1: 启动服务
```bash
# 在终端中运行
sudo ./build/clash-fish start
```

**预期行为**:
1. 显示启动信息
2. 检测并显示 VPN 状态
3. 启动 mihomo 引擎
4. 创建 TUN 设备（utun5）
5. 配置路由表
6. 显示 "Service is running in foreground"
7. 等待 Ctrl+C 信号

**验证点**:
- [ ] VPN 检测信息正确显示
- [ ] 没有错误信息
- [ ] PID 文件创建成功
- [ ] 服务在前台运行

#### 测试 2: 检查运行状态（另开终端）
```bash
# 在另一个终端运行
./build/clash-fish status
```

**预期输出**:
```
Service:    ✓ Running (PID: xxxxx)
VPN:        ✓ Active (utun4: 10.8.0.105)
Mode:       rule
HTTP Port:  7890
SOCKS Port: 7891
TUN Mode:   true
```

#### 测试 3: 验证 TUN 设备
```bash
ifconfig | grep utun
```

**预期**:
- 应该看到 utun5 接口
- IP 地址为 198.18.0.1

#### 测试 4: 验证路由表
```bash
netstat -nr | grep 198.18
```

**预期**:
- 应该看到指向 utun5 的路由规则
- 0.0.0.0/1 和 128.0.0.0/1 应该指向 198.18.0.1

#### 测试 5: 验证代理端口
```bash
lsof -i :7890
lsof -i :7891
```

**预期**:
- 端口 7890 (HTTP) 被占用
- 端口 7891 (SOCKS5) 被占用

#### 测试 6: 测试代理功能
```bash
# 设置代理环境变量
export http_proxy=http://127.0.0.1:7890
export https_proxy=http://127.0.0.1:7890

# 测试连接
curl -I https://www.google.com
```

**预期**:
- 能够成功访问 Google
- 流量通过代理

#### 测试 7: 验证 VPN 共存
```bash
# 测试访问 VPN 内网地址（如果有的话）
ping 10.8.0.1

# 同时测试外网
curl https://www.google.com
```

**预期**:
- VPN 内网地址可访问（走 utun4）
- 外网地址也可访问（走 utun5 代理）

#### 测试 8: 停止服务
```bash
# 方式 1: 在运行终端按 Ctrl+C

# 方式 2: 在另一个终端运行
sudo ./build/clash-fish stop
```

**预期**:
- 显示 "Stopping Clash-Fish..."
- 优雅关闭
- PID 文件被删除
- 显示 "✓ Clash-Fish stopped"

#### 测试 9: 验证清理
```bash
./build/clash-fish status
ls ~/.config/clash-fish/*.pid
```

**预期**:
- 状态显示 "Not Running"
- PID 文件不存在

#### 测试 10: 重启服务
```bash
sudo ./build/clash-fish restart
```

**预期**:
- 先停止（如果在运行）
- 然后启动
- 显示成功信息

---

## 已知限制

1. **需要 sudo 权限**
   - TUN 模式需要 root 权限创建虚拟网卡
   - 路由表修改需要 root 权限

2. **前台运行**
   - 当前版本在前台运行
   - 需要 Ctrl+C 或另一个终端执行 stop 命令
   - 后续可以添加 daemon 模式

3. **日志输出**
   - mihomo 日志当前输出到终端
   - 后续可以重定向到文件

---

## 测试检查清单

### 基础功能
- [x] 配置初始化
- [x] 配置验证
- [x] 配置显示
- [x] VPN 检测
- [x] 权限检查
- [ ] 服务启动（需要 sudo）
- [ ] 服务停止（需要 sudo）
- [ ] 服务重启（需要 sudo）
- [ ] 运行状态检查

### 网络功能
- [ ] TUN 设备创建
- [ ] 路由表配置
- [ ] HTTP 代理 (7890)
- [ ] SOCKS5 代理 (7891)
- [ ] DNS 劫持
- [ ] 代理连接测试

### VPN 共存
- [x] VPN 检测功能
- [ ] VPN 流量正确路由
- [ ] 代理流量正确路由
- [ ] 同时访问内外网

---

## 下一步建议

### 优先级 P0（立即）
1. 用户手动测试 `sudo clash-fish start`
2. 验证 TUN 设备和路由表
3. 测试代理连接

### 优先级 P1（Day 3）
1. 实现订阅导入功能
2. 实现 Profile 切换
3. 添加更多错误处理

### 优先级 P2（后续）
1. 添加 daemon 模式（后台运行）
2. 添加日志文件输出
3. 添加流量统计
4. 添加系统服务安装（launchd）

---

## 技术细节

### 配置文件结构
```
~/.config/clash-fish/
├── config.yaml           # 主配置文件
├── logs/
│   ├── clash-fish.log    # 应用日志
│   └── mihomo.log        # mihomo 日志（如果配置）
├── profiles/             # 多配置文件
├── cache/                # 缓存
└── clash-fish.pid        # PID 文件（运行时）
```

### 预期的网络配置

**TUN 设备**:
```
utun5: flags=8051<UP,POINTOPOINT,RUNNING,MULTICAST>
    inet 198.18.0.1 --> 198.18.0.1 netmask 0xffff0000
```

**路由规则**:
```
# VPN 路由（高优先级）
10.8.0.0/24        10.8.0.x       UGSc      utun4
192.168.0.0/16     10.8.0.1       UGSc      utun4

# 代理路由（次优先级）
0.0.0.0/1          198.18.0.1     UGSc      utun5
128.0.0.0/1        198.18.0.1     UGSc      utun5

# 默认路由（低优先级）
default            192.168.1.1    UGScg     en0
```

**DNS 配置**:
```
nameserver[0] : 198.18.0.2    # clash-fish DNS
```

---

## 故障排查

### 问题: 启动失败

**检查**:
1. 是否使用 sudo
2. 配置文件是否存在
3. 配置是否有效
4. 端口是否被占用

**解决**:
```bash
# 检查端口占用
lsof -i :7890
lsof -i :7891

# 验证配置
./build/clash-fish config validate

# 查看详细日志
sudo ./build/clash-fish start --debug
```

### 问题: VPN 流量被代理劫持

**检查**:
```bash
netstat -nr
```

**解决**:
- 确保 VPN 路由优先级更高
- 检查 VPN 路由是否正确配置

### 问题: 无法访问外网

**检查**:
1. TUN 设备是否创建
2. 路由表是否正确
3. 代理服务器是否可用

---

**测试人员**: Claude Code
**测试日期**: 2025-11-23
**测试状态**: 部分通过（需要 sudo 权限完成完整测试）
