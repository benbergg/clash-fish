# Clash-Fish 技术设计文档

## 1. 项目概述

### 1.1 项目目标
开发一个基于 mihomo 内核的 macOS 透明代理工具，实现 VPN 和透明代理共存，支持自动导入代理配置，提供命令行界面。

### 1.2 核心特性
- ✅ macOS 原生支持
- ✅ 基于 mihomo (Clash Meta) 内核
- ✅ VPN 和透明代理共存
- ✅ 支持 nolock 代理模式
- ✅ 命令行界面（CLI）
- ✅ 自动导入代理配置
- ✅ 快速开发周期

### 1.3 项目命名
**Clash-Fish** - 寓意灵活、快速的代理工具

---

## 2. 技术架构

### 2.1 整体架构

```
┌─────────────────────────────────────────────────────────┐
│                    CLI Interface                         │
│          (配置管理、启动/停止、状态查看)                   │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                  Core Service Layer                      │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │ Config Mgr   │  │ Proxy Mgr    │  │ System Mgr   │  │
│  │ (配置管理)   │  │ (代理管理)   │  │ (系统设置)   │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│                   Mihomo Core                            │
│    (Clash Meta 内核 - 代理转发、规则匹配)                │
└─────────────────────────────────────────────────────────┘
                          │
                          ▼
┌─────────────────────────────────────────────────────────┐
│              macOS System Integration                    │
│  ┌──────────────┐  ┌──────────────┐  ┌──────────────┐  │
│  │   TUN Mode   │  │  System      │  │  DNS         │  │
│  │   (虚拟网卡) │  │  Proxy       │  │  Hijack      │  │
│  └──────────────┘  └──────────────┘  └──────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### 2.2 技术栈

| 层次 | 技术选型 | 说明 |
|------|---------|------|
| 编程语言 | Go 1.21+ | 性能好，mihomo 原生支持 |
| 核心引擎 | Mihomo Core | Clash Meta 分支，功能强大 |
| CLI 框架 | cobra + viper | 标准 Go CLI 工具链 |
| 配置格式 | YAML | 兼容 Clash 配置格式 |
| 系统集成 | TUN/utun | macOS 虚拟网卡 |
| 权限管理 | sudo/osascript | macOS 提权机制 |
| 日志 | zerolog | 高性能结构化日志 |

---

## 3. 核心功能模块

### 3.1 配置管理模块

#### 功能
- 读取、验证、更新 mihomo 配置文件
- 支持多配置文件管理（profile 切换）
- 自动导入订阅链接（支持 Clash/V2Ray 订阅）
- 配置模板生成

#### 实现方案
```go
type ConfigManager struct {
    configPath   string
    profiles     map[string]*Profile
    activeProfile string
}

type Profile struct {
    Name     string
    Path     string
    URL      string  // 订阅链接
    UpdatedAt time.Time
}

// 核心方法
- LoadConfig(path string) error
- ValidateConfig(config *Config) error
- UpdateSubscription(url string) error
- SwitchProfile(name string) error
- GenerateTemplate() error
```

### 3.2 代理管理模块

#### 功能
- mihomo 进程生命周期管理
- 代理状态监控
- 流量统计
- 连接管理

#### 实现方案
```go
type ProxyManager struct {
    core      *mihomo.Engine
    status    ProxyStatus
    statsCollector *StatsCollector
}

// 核心方法
- Start() error
- Stop() error
- Restart() error
- GetStatus() ProxyStatus
- GetTraffic() TrafficStats
```

### 3.3 系统集成模块

#### 功能
- TUN 模式实现（透明代理）
- 系统代理设置（HTTP/HTTPS/SOCKS5）
- DNS 劫持配置
- 路由表管理

#### 实现方案

**TUN 模式实现**（透明代理核心）
```go
type TUNManager struct {
    device    string  // utun 设备名
    tunFd     int
    routes    []Route
}

// 核心方法
- CreateTUN() error              // 创建 TUN 设备
- ConfigureRoutes() error        // 配置路由规则
- SetupIPTables() error          // 配置流量转发规则
- Cleanup() error                // 清理系统设置
```

**系统代理设置**
```go
type SystemProxyManager struct {
    networkService string
}

// 使用 networksetup 命令
- SetHTTPProxy(host, port string) error
- SetHTTPSProxy(host, port string) error
- SetSOCKSProxy(host, port string) error
- DisableProxy() error
```

### 3.4 Nolock 代理模式

#### 说明
Nolock 模式指的是 TCP no-delay 模式，减少延迟，适合游戏等实时应用。

#### 实现方案
```yaml
# mihomo 配置示例
proxies:
  - name: "nolock-proxy"
    type: ss
    server: server
    port: 443
    cipher: aes-256-gcm
    password: password
    udp: true
    plugin: v2ray-plugin
    plugin-opts:
      mode: websocket
      mux: false  # 关键：禁用多路复用

tcp-concurrent: true  # mihomo 支持的 TCP 并发
```

---

## 4. VPN 与透明代理共存方案

### 4.1 技术挑战
macOS 系统同时运行 VPN 和透明代理会导致路由冲突。

### 4.2 方案选择

**基于实际环境测试，clash-fish 采用：方案一 - 基于路由优先级**

#### 选择理由
1. ✅ 已在 macOS 15.7.2 环境验证可行
2. ✅ 实现简单，mihomo 原生支持
3. ✅ 性能最优（内核层面路由决策）
4. ✅ 完全自动化，无需手动配置
5. ✅ 符合"内网 VPN + 外网代理"的典型场景

### 4.3 工作原理

#### 核心机制：利用 IP 路由表的"最长前缀匹配"规则

```
路由优先级（掩码越长优先级越高）：
┌────────────────────────────────────────────────┐
│ VPN 内网路由 (最高优先级)                        │
│ 10.8.0.0/24      → utun4 (VPN)                 │
│ 192.168.0.0/16   → utun4 (VPN)                 │
├────────────────────────────────────────────────┤
│ 代理路由 (次高优先级)                            │
│ 0.0.0.0/1        → utun5 (Proxy)               │
│ 128.0.0.0/1      → utun5 (Proxy)               │
├────────────────────────────────────────────────┤
│ 默认路由 (最低优先级)                            │
│ default          → en0 (Gateway)                │
└────────────────────────────────────────────────┘

流量走向示例：
访问 10.8.0.50        → VPN (utun4)     [内网服务器]
访问 192.168.0.100    → VPN (utun4)     [内网资源]
访问 8.8.8.8          → PROXY (utun5)   [公网 DNS]
访问 google.com       → PROXY (utun5)   [代理]
访问 192.168.1.1      → DIRECT (en0)    [本地网关]
```

### 4.4 实现方案

#### mihomo 配置
```yaml
tun:
  enable: true
  stack: system               # 使用系统网络栈
  auto-route: true            # 自动配置路由（关键！）
  auto-detect-interface: true # 自动检测网卡
  dns-hijack:
    - any:53

dns:
  enable: true
  listen: 198.18.0.2:53
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16

rules:
  - GEOIP,PRIVATE,DIRECT      # 私网地址直连（包括 VPN 内网）
  - GEOIP,CN,DIRECT           # 可选：国内直连
  - MATCH,PROXY
```

#### 代码实现
```go
// VPN 检测
func DetectVPN() (bool, string) {
    interfaces, _ := net.Interfaces()
    for _, iface := range interfaces {
        if strings.HasPrefix(iface.Name, "utun") {
            addrs, _ := iface.Addrs()
            for _, addr := range addrs {
                ip := parseIP(addr)
                // 检测是否是私网地址（VPN 特征）
                if isPrivateIP(ip) {
                    return true, iface.Name
                }
            }
        }
    }
    return false, ""
}

// 启动时智能提示
func Start() error {
    vpnActive, vpnDev := DetectVPN()
    if vpnActive {
        log.Info().
            Str("device", vpnDev).
            Msg("检测到 VPN 连接，将自动与代理共存")
    }

    // mihomo 的 auto-route 会自动处理路由配置
    startMihomo()
}
```

### 4.5 兼容性

已测试兼容的 VPN 类型：
- ✅ OpenVPN
- ✅ WireGuard
- ✅ Cisco AnyConnect
- ✅ L2TP/IPSec

### 4.6 未来扩展（可选）

如果用户有复杂的域名级分流需求，可在后续版本添加**规则分流模式**：

```yaml
# 高级配置（v2.0 考虑）
vpn-coexist-mode: routing  # routing(默认) | rule

# 规则模式配置
vpn-rules:
  - DOMAIN-SUFFIX,company.com,VPN
  - IP-CIDR,10.0.0.0/8,VPN
```

---

## 5. 自动导入代理配置

### 5.1 支持的订阅格式
- Clash 订阅（YAML）
- V2Ray 订阅（base64 编码）
- Shadowsocks 订阅（SIP002）
- Trojan 订阅

### 5.2 实现流程
```go
func ImportSubscription(url string) error {
    // 1. 下载订阅内容
    content := httpClient.Get(url)

    // 2. 检测格式
    format := DetectFormat(content)

    // 3. 转换为 mihomo 格式
    config := ConvertToMihomo(content, format)

    // 4. 验证配置
    ValidateConfig(config)

    // 5. 保存配置文件
    SaveConfig(config)

    // 6. 重载服务
    ReloadProxy()
}
```

### 5.3 定时更新
```go
type SubscriptionUpdater struct {
    interval time.Duration
    ticker   *time.Ticker
}

// 每 24 小时自动更新订阅
func (u *SubscriptionUpdater) Start() {
    u.ticker = time.NewTicker(24 * time.Hour)
    go func() {
        for range u.ticker.C {
            UpdateAllSubscriptions()
        }
    }()
}
```

---

## 6. CLI 命令设计

### 6.1 命令结构
```bash
clash-fish                          # 主命令
├── start                           # 启动服务
├── stop                            # 停止服务
├── restart                         # 重启服务
├── status                          # 查看状态
├── config                          # 配置管理
│   ├── init                        # 初始化配置
│   ├── edit                        # 编辑配置
│   ├── validate                    # 验证配置
│   └── show                        # 显示配置
├── profile                         # 配置文件管理
│   ├── list                        # 列出所有 profile
│   ├── add <name> <url>           # 添加订阅
│   ├── update [name]              # 更新订阅
│   ├── switch <name>              # 切换 profile
│   └── delete <name>              # 删除 profile
├── proxy                           # 代理管理
│   ├── list                        # 列出所有代理
│   ├── test [name]                # 测试延迟
│   └── select <group> <proxy>     # 选择代理
├── mode                            # 代理模式
│   ├── set <mode>                 # 设置模式 (global/rule/direct)
│   └── show                       # 显示当前模式
└── log                            # 日志查看
    ├── tail                       # 实时日志
    └── show                       # 历史日志
```

### 6.2 使用示例
```bash
# 初始化配置
clash-fish config init

# 添加订阅
clash-fish profile add my-sub https://example.com/clash

# 启动服务
sudo clash-fish start

# 查看状态
clash-fish status

# 切换规则模式
clash-fish mode set rule

# 查看实时日志
clash-fish log tail
```

---

## 7. 项目结构

```
clash-fish/
├── cmd/
│   └── clash-fish/
│       ├── main.go              # 主入口
│       ├── start.go             # start 命令
│       ├── stop.go              # stop 命令
│       ├── config.go            # config 命令组
│       └── profile.go           # profile 命令组
├── internal/
│   ├── config/
│   │   ├── manager.go           # 配置管理器
│   │   ├── parser.go            # 配置解析
│   │   └── validator.go         # 配置验证
│   ├── proxy/
│   │   ├── manager.go           # 代理管理器
│   │   ├── mihomo.go            # mihomo 引擎封装
│   │   └── stats.go             # 流量统计
│   ├── system/
│   │   ├── tun.go               # TUN 设备管理
│   │   ├── route.go             # 路由管理
│   │   ├── dns.go               # DNS 配置
│   │   └── sysproxy.go          # 系统代理设置
│   ├── subscription/
│   │   ├── importer.go          # 订阅导入
│   │   ├── converter.go         # 格式转换
│   │   └── updater.go           # 定时更新
│   └── service/
│       ├── daemon.go            # 后台服务
│       └── controller.go        # 服务控制器
├── pkg/
│   ├── logger/                  # 日志工具
│   ├── utils/                   # 通用工具
│   └── constants/               # 常量定义
├── configs/
│   ├── config.example.yaml      # 配置模板
│   └── rules/                   # 规则文件
├── scripts/
│   ├── install.sh               # 安装脚本
│   └── uninstall.sh             # 卸载脚本
├── go.mod
├── go.sum
├── Makefile
└── README.md
```

---

## 8. 依赖管理

### 8.1 核心依赖
```go
require (
    github.com/metacubex/mihomo v1.18.0  // mihomo 核心
    github.com/spf13/cobra v1.8.0        // CLI 框架
    github.com/spf13/viper v1.18.0       // 配置管理
    github.com/rs/zerolog v1.32.0        // 日志库
    gopkg.in/yaml.v3 v3.0.1              // YAML 解析
)
```

### 8.2 mihomo 集成方式

**方式一：作为库引入（推荐）**
```go
import (
    "github.com/metacubex/mihomo/config"
    "github.com/metacubex/mihomo/hub/executor"
    "github.com/metacubex/mihomo/log"
)

func StartMihomo(configPath string) error {
    cfg, err := config.Parse(configPath)
    if err != nil {
        return err
    }
    executor.ApplyConfig(cfg, true)
    return nil
}
```

**方式二：独立进程（备选）**
```go
// 启动 mihomo 二进制
cmd := exec.Command("mihomo", "-d", configDir, "-f", configFile)
cmd.Start()
```

---

## 9. 开发计划

### 9.1 开发策略

**采用敏捷开发，优先交付 MVP 版本**

- **MVP 目标**：3-5 天完成基础可用版本
- **完整版本**：8-10 天完成所有核心功能
- **迭代优化**：2-3 天进行测试和优化

### 9.2 里程碑规划

#### Phase 1: 项目初始化（Day 1 上午）
- [x] 项目目录创建
- [ ] Go 模块初始化
- [ ] 依赖包安装（cobra, viper, zerolog, mihomo）
- [ ] 基础项目结构搭建
- [ ] Git 仓库初始化

**交付物**: 可编译的空壳项目

#### Phase 2: CLI 框架（Day 1 下午）
- [ ] Cobra CLI 命令结构
  - [ ] start 命令
  - [ ] stop 命令
  - [ ] status 命令
  - [ ] config 命令组
- [ ] Viper 配置加载
- [ ] Zerolog 日志系统
- [ ] 基础常量定义

**交付物**: 可运行的 CLI 框架（命令结构完整但功能为空）

#### Phase 3: 配置管理（Day 2 上午）
- [ ] 配置文件模板生成
- [ ] YAML 解析器
- [ ] 配置验证器
- [ ] 配置文件读写
- [ ] 配置目录管理 (~/.config/clash-fish)

**交付物**: `clash-fish config init` 可用

#### Phase 4: Mihomo 集成（Day 2 下午 - Day 3）
- [ ] Mihomo 引擎封装
- [ ] TUN 模式启动
- [ ] 进程管理（启动/停止/重启）
- [ ] PID 文件管理
- [ ] 权限检查（root）

**交付物**: `clash-fish start/stop` 基础功能可用

#### Phase 5: 系统集成（Day 3 下午 - Day 4 上午）
- [ ] VPN 检测功能
- [ ] 路由表查看工具
- [ ] DNS 配置验证
- [ ] 优雅关闭（信号处理）
- [ ] 清理机制（退出时恢复设置）

**交付物**: VPN 共存功能验证通过

#### Phase 6: 订阅管理（Day 4 下午 - Day 5 上午）
- [ ] HTTP 订阅下载
- [ ] Clash 格式解析
- [ ] 配置文件保存
- [ ] Profile 管理（add/list/switch/delete）
- [ ] `clash-fish profile add` 命令

**交付物**: 订阅导入功能可用

#### Phase 7: 状态与监控（Day 5 下午）
- [ ] 服务状态检查
- [ ] 流量统计（基础）
- [ ] 日志查看命令
- [ ] 错误处理完善

**交付物**: `clash-fish status` 和 `clash-fish log` 可用

#### Phase 8: 测试与优化（Day 6-7）
- [ ] VPN 共存场景测试
- [ ] 订阅导入测试
- [ ] 错误场景测试
- [ ] 性能测试
- [ ] 文档完善

**交付物**: MVP 1.0 版本发布

### 9.3 MVP 功能清单

**必须有（P0）**:
1. ✅ 启动/停止 mihomo 核心
2. ✅ TUN 模式透明代理
3. ✅ 基本配置文件管理
4. ✅ CLI 命令：start/stop/status/config init
5. ✅ VPN 自动检测和共存
6. ✅ 订阅导入（Clash 格式）

**应该有（P1）**:
1. Profile 切换功能
2. 日志查看功能
3. 权限自动提升
4. 配置验证功能

**可以有（P2）**:
1. 流量统计
2. 代理延迟测试
3. 规则模式切换
4. 定时订阅更新

**暂不做**:
1. GUI 界面
2. V2Ray/SS 订阅格式支持
3. 高级规则编辑
4. 性能优化

### 9.4 每日开发任务

详见 `DEVELOPMENT_PLAN.md`

---

## 10. 技术难点与解决方案

### 10.1 macOS 权限问题

**问题**：TUN 设备创建需要 root 权限

**解决方案**：
```go
// 检测权限
func CheckRoot() bool {
    return os.Geteuid() == 0
}

// 提权执行
func ElevatePrivileges() error {
    if !CheckRoot() {
        cmd := exec.Command("osascript", "-e",
            `do shell script "sudo /path/to/clash-fish start" with administrator privileges`)
        return cmd.Run()
    }
    return nil
}
```

### 10.2 TUN 设备配置

**问题**：macOS 的 utun 设备配置复杂

**解决方案**：
```go
// 使用 mihomo 内置的 TUN 实现
type TUNConfig struct {
    Enable     bool     `yaml:"enable"`
    Stack      string   `yaml:"stack"`      // system/gvisor
    DNSHijack  []string `yaml:"dns-hijack"`
    AutoRoute  bool     `yaml:"auto-route"`
    AutoDetectInterface bool `yaml:"auto-detect-interface"`
}

// mihomo 会自动处理 TUN 设备的创建和配置
```

### 10.3 DNS 劫持

**问题**：DNS 查询需要被代理

**解决方案**：
```yaml
dns:
  enable: true
  listen: 0.0.0.0:53
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  nameserver:
    - 223.5.5.5
    - 114.114.114.114
  fallback:
    - 8.8.8.8
    - 1.1.1.1
```

### 10.4 与 VPN 共存

**问题**：路由表冲突

**解决方案**：
```go
// 检测 VPN 连接
func DetectVPN() bool {
    // 检查是否有 ppp/ipsec 接口
    interfaces, _ := net.Interfaces()
    for _, iface := range interfaces {
        if strings.HasPrefix(iface.Name, "ppp") ||
           strings.HasPrefix(iface.Name, "ipsec") {
            return true
        }
    }
    return false
}

// 动态调整路由策略
func AdjustRoutingPolicy(vpnActive bool) {
    if vpnActive {
        // 使用基于规则的路由，不劫持全部流量
        useTUNMode = false
        useSystemProxy = true
    } else {
        // VPN 未连接，可以使用 TUN 全局代理
        useTUNMode = true
    }
}
```

---

## 11. 安全考虑

### 11.1 配置文件安全
- 配置文件权限：`chmod 600 config.yaml`
- 敏感信息加密存储（密码、密钥）
- 避免明文记录代理密码

### 11.2 进程安全
- 以最小权限运行（仅必要操作提权）
- 进程 PID 文件防止重复启动
- 优雅关闭处理（信号捕获）

### 11.3 网络安全
- DNS 泄漏防护
- WebRTC IP 泄漏防护
- 流量加密验证

---

## 12. 性能优化

### 12.1 内存优化
- 使用 mihomo 的流式处理
- 限制连接数和缓冲区大小
- 定期清理过期连接

### 12.2 延迟优化
- Nolock 模式：禁用 Nagle 算法
- TCP Fast Open
- 连接池复用

### 12.3 并发优化
```yaml
# mihomo 配置
tcp-concurrent: true
sniffer:
  enable: true
  sniffing:
    - tls
    - http
  port-whitelist:
    - 80
    - 443
```

---

## 13. 测试计划

### 13.1 单元测试
- 配置解析测试
- 订阅转换测试
- 路由规则测试

### 13.2 集成测试
- mihomo 引擎启动测试
- TUN 设备创建测试
- 代理连接测试

### 13.3 兼容性测试
- macOS 12 (Monterey)
- macOS 13 (Ventura)
- macOS 14 (Sonoma)
- macOS 15 (Sequoia)

### 13.4 VPN 共存测试
- + Cisco AnyConnect
- + OpenVPN
- + WireGuard
- + L2TP/IPSec

---

## 14. 部署与分发

### 14.1 编译
```makefile
# Makefile
build:
	go build -o clash-fish cmd/clash-fish/main.go

install:
	sudo cp clash-fish /usr/local/bin/
	sudo chmod +x /usr/local/bin/clash-fish

uninstall:
	sudo rm /usr/local/bin/clash-fish
```

### 14.2 安装脚本
```bash
#!/bin/bash
# install.sh

# 1. 下载二进制
curl -L https://github.com/user/clash-fish/releases/latest/download/clash-fish -o clash-fish

# 2. 安装到系统
sudo mv clash-fish /usr/local/bin/
sudo chmod +x /usr/local/bin/clash-fish

# 3. 创建配置目录
mkdir -p ~/.config/clash-fish

# 4. 初始化配置
clash-fish config init

echo "Installation complete!"
```

### 14.3 后续 GUI 版本考虑
- SwiftUI macOS 应用
- 系统托盘图标
- 菜单栏快捷操作
- 图形化规则编辑器

---

## 15. 维护与支持

### 15.1 日志位置
```
~/.config/clash-fish/logs/
├── clash-fish.log          # 应用日志
├── mihomo.log              # mihomo 核心日志
└── error.log               # 错误日志
```

### 15.2 配置位置
```
~/.config/clash-fish/
├── config.yaml             # 主配置
├── profiles/               # 多配置文件
│   ├── default.yaml
│   └── work.yaml
└── cache/                  # 缓存文件
```

### 15.3 故障排查
```bash
# 查看服务状态
clash-fish status

# 查看日志
clash-fish log tail

# 验证配置
clash-fish config validate

# 测试代理连接
clash-fish proxy test
```

---

## 16. 参考资源

### 16.1 相关项目
- [mihomo](https://github.com/MetaCubeX/mihomo) - Clash Meta 内核
- [ClashX](https://github.com/yichengchen/clashX) - macOS GUI 参考
- [Clash Verge](https://github.com/zzzgydi/clash-verge) - 跨平台参考

### 16.2 技术文档
- [Clash 配置文档](https://clash.wiki/)
- [mihomo Wiki](https://wiki.metacubex.one/)
- [macOS Network Extension](https://developer.apple.com/documentation/networkextension)

### 16.3 协议规范
- [SOCKS5 RFC](https://datatracker.ietf.org/doc/html/rfc1928)
- [Shadowsocks SIP002](https://shadowsocks.org/doc/sip002.html)
- [V2Ray Protocol](https://www.v2ray.com/developer/protocols/)

---

## 附录 A: 配置文件示例

```yaml
# config.yaml
port: 7890
socks-port: 7891
allow-lan: false
mode: rule
log-level: info
external-controller: 127.0.0.1:9090

# TUN 配置
tun:
  enable: true
  stack: system
  dns-hijack:
    - any:53
  auto-route: true
  auto-detect-interface: true

# DNS 配置
dns:
  enable: true
  listen: 0.0.0.0:53
  enhanced-mode: fake-ip
  fake-ip-range: 198.18.0.1/16
  nameserver:
    - 223.5.5.5
    - 114.114.114.114
  fallback:
    - tls://1.1.1.1:853
    - tls://8.8.8.8:853

# 代理配置
proxies:
  - name: "ss-nolock"
    type: ss
    server: example.com
    port: 8388
    cipher: aes-256-gcm
    password: password
    udp: true

# 代理组
proxy-groups:
  - name: "PROXY"
    type: select
    proxies:
      - "ss-nolock"
      - DIRECT

# 规则
rules:
  - DOMAIN-SUFFIX,google.com,PROXY
  - DOMAIN-KEYWORD,youtube,PROXY
  - GEOIP,CN,DIRECT
  - MATCH,PROXY
```

---

## 附录 B: 快速开始指南

```bash
# 1. 安装
curl -L https://github.com/user/clash-fish/releases/latest/download/install.sh | bash

# 2. 初始化配置
clash-fish config init

# 3. 导入订阅
clash-fish profile add my-sub https://example.com/clash

# 4. 启动服务（需要 sudo）
sudo clash-fish start

# 5. 查看状态
clash-fish status

# 6. 设置代理模式
clash-fish mode set rule

# 7. 查看代理列表
clash-fish proxy list

# 8. 停止服务
sudo clash-fish stop
```

---

**文档版本**: v1.0
**创建时间**: 2025-11-23
**最后更新**: 2025-11-23
**维护者**: clash-fish team
