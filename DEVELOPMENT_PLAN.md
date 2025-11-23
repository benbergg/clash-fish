# Clash-Fish 开发执行计划

## 项目信息
- **项目名称**: Clash-Fish
- **开发周期**: 7-8 天
- **MVP 目标**: Day 5 完成
- **测试优化**: Day 6-7
- **开发模式**: 敏捷迭代

---

## Day 1: 项目初始化与 CLI 框架

### 上午任务：项目初始化（2-3小时）

#### 1.1 创建项目结构
```bash
cd /Users/lg/Projects/go/clash-fish

# 初始化 Go 模块
go mod init github.com/yourusername/clash-fish

# 创建目录结构
mkdir -p cmd/clash-fish
mkdir -p internal/{config,proxy,system,subscription,service}
mkdir -p pkg/{logger,utils,constants}
mkdir -p configs/rules
mkdir -p scripts
```

#### 1.2 安装依赖
```bash
# 核心依赖
go get github.com/metacubex/mihomo@latest
go get github.com/spf13/cobra@latest
go get github.com/spf13/viper@latest
go get github.com/rs/zerolog@latest
go get gopkg.in/yaml.v3@latest

# 工具依赖
go get github.com/google/uuid@latest
```

#### 1.3 创建基础文件
- [ ] `cmd/clash-fish/main.go` - 主入口
- [ ] `pkg/logger/logger.go` - 日志工具
- [ ] `pkg/constants/constants.go` - 常量定义
- [ ] `Makefile` - 构建脚本
- [ ] `.gitignore` - Git 忽略文件

**检查点**: `go build` 成功，程序可运行

---

### 下午任务：CLI 框架（3-4小时）

#### 1.4 实现 Cobra CLI 结构

**文件**: `cmd/clash-fish/main.go`
```go
package main

import (
    "fmt"
    "os"
    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "clash-fish",
    Short: "A transparent proxy tool based on mihomo",
    Long:  `Clash-Fish is a CLI tool for macOS that provides transparent proxy with VPN coexistence.`,
}

func main() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Fprintln(os.Stderr, err)
        os.Exit(1)
    }
}
```

#### 1.5 实现子命令

- [ ] **start.go** - 启动服务
- [ ] **stop.go** - 停止服务
- [ ] **status.go** - 查看状态
- [ ] **restart.go** - 重启服务
- [ ] **config.go** - 配置管理命令组
- [ ] **profile.go** - Profile 管理命令组
- [ ] **log.go** - 日志查看

#### 1.6 日志系统

**文件**: `pkg/logger/logger.go`
```go
package logger

import (
    "os"
    "path/filepath"
    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func Init(logDir string) error {
    // 创建日志目录
    if err := os.MkdirAll(logDir, 0755); err != nil {
        return err
    }

    // 配置日志
    logFile := filepath.Join(logDir, "clash-fish.log")
    file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        return err
    }

    // 多输出：控制台 + 文件
    multi := zerolog.MultiLevelWriter(
        zerolog.ConsoleWriter{Out: os.Stdout},
        file,
    )

    log.Logger = zerolog.New(multi).With().Timestamp().Caller().Logger()
    zerolog.SetGlobalLevel(zerolog.InfoLevel)

    return nil
}
```

**检查点**:
- `clash-fish --help` 显示帮助信息
- `clash-fish start --help` 显示 start 命令帮助
- 日志系统可用

**Day 1 交付物**: CLI 框架完整，所有命令有框架但功能为空

---

## Day 2: 配置管理与 Mihomo 集成

### 上午任务：配置管理（3小时）

#### 2.1 配置结构定义

**文件**: `internal/config/types.go`
```go
package config

type Config struct {
    Port         int           `yaml:"port"`
    SocksPort    int           `yaml:"socks-port"`
    AllowLan     bool          `yaml:"allow-lan"`
    Mode         string        `yaml:"mode"`
    LogLevel     string        `yaml:"log-level"`
    ExternalController string `yaml:"external-controller"`
    TUN          TUNConfig     `yaml:"tun"`
    DNS          DNSConfig     `yaml:"dns"`
    Proxies      []Proxy       `yaml:"proxies"`
    ProxyGroups  []ProxyGroup  `yaml:"proxy-groups"`
    Rules        []string      `yaml:"rules"`
}

type TUNConfig struct {
    Enable              bool     `yaml:"enable"`
    Stack               string   `yaml:"stack"`
    DNSHijack           []string `yaml:"dns-hijack"`
    AutoRoute           bool     `yaml:"auto-route"`
    AutoDetectInterface bool     `yaml:"auto-detect-interface"`
}

type DNSConfig struct {
    Enable        bool     `yaml:"enable"`
    Listen        string   `yaml:"listen"`
    EnhancedMode  string   `yaml:"enhanced-mode"`
    FakeIPRange   string   `yaml:"fake-ip-range"`
    Nameserver    []string `yaml:"nameserver"`
    Fallback      []string `yaml:"fallback"`
}
```

#### 2.2 配置管理器

**文件**: `internal/config/manager.go`
```go
package config

import (
    "os"
    "path/filepath"
    "gopkg.in/yaml.v3"
)

type Manager struct {
    configDir  string
    configPath string
    config     *Config
}

func NewManager(configDir string) *Manager {
    return &Manager{
        configDir:  configDir,
        configPath: filepath.Join(configDir, "config.yaml"),
    }
}

func (m *Manager) Init() error {
    // 创建配置目录
    if err := os.MkdirAll(m.configDir, 0755); err != nil {
        return err
    }

    // 如果配置文件不存在，创建默认配置
    if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
        return m.CreateDefault()
    }

    return nil
}

func (m *Manager) CreateDefault() error {
    config := GetDefaultConfig()
    return m.Save(config)
}

func (m *Manager) Load() (*Config, error) {
    data, err := os.ReadFile(m.configPath)
    if err != nil {
        return nil, err
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }

    m.config = &config
    return &config, nil
}

func (m *Manager) Save(config *Config) error {
    data, err := yaml.Marshal(config)
    if err != nil {
        return err
    }

    return os.WriteFile(m.configPath, data, 0644)
}

func (m *Manager) Validate(config *Config) error {
    // TODO: 实现配置验证逻辑
    return nil
}
```

#### 2.3 默认配置模板

**文件**: `internal/config/template.go`
```go
package config

func GetDefaultConfig() *Config {
    return &Config{
        Port:               7890,
        SocksPort:          7891,
        AllowLan:           false,
        Mode:               "rule",
        LogLevel:           "info",
        ExternalController: "127.0.0.1:9090",
        TUN: TUNConfig{
            Enable:              true,
            Stack:               "system",
            DNSHijack:           []string{"any:53"},
            AutoRoute:           true,
            AutoDetectInterface: true,
        },
        DNS: DNSConfig{
            Enable:       true,
            Listen:       "198.18.0.2:53",
            EnhancedMode: "fake-ip",
            FakeIPRange:  "198.18.0.1/16",
            Nameserver:   []string{"223.5.5.5", "114.114.114.114"},
            Fallback:     []string{"tls://1.1.1.1:853", "tls://8.8.8.8:853"},
        },
        ProxyGroups: []ProxyGroup{
            {
                Name:    "PROXY",
                Type:    "select",
                Proxies: []string{"DIRECT"},
            },
        },
        Rules: []string{
            "GEOIP,PRIVATE,DIRECT",
            "GEOIP,CN,DIRECT",
            "MATCH,PROXY",
        },
    }
}
```

**检查点**: `clash-fish config init` 创建默认配置文件

---

### 下午任务：Mihomo 集成（4小时）

#### 2.4 Mihomo 引擎封装

**文件**: `internal/proxy/mihomo.go`
```go
package proxy

import (
    "github.com/metacubex/mihomo/config"
    "github.com/metacubex/mihomo/hub/executor"
    "github.com/metacubex/mihomo/log"
)

type MihomoEngine struct {
    configPath string
    running    bool
}

func NewMihomoEngine(configPath string) *MihomoEngine {
    return &MihomoEngine{
        configPath: configPath,
        running:    false,
    }
}

func (e *MihomoEngine) Start() error {
    // 解析配置
    cfg, err := config.Parse(e.configPath)
    if err != nil {
        return err
    }

    // 应用配置
    executor.ApplyConfig(cfg, true)

    e.running = true
    log.Infoln("Mihomo engine started")

    return nil
}

func (e *MihomoEngine) Stop() error {
    // TODO: 实现停止逻辑
    e.running = false
    return nil
}

func (e *MihomoEngine) IsRunning() bool {
    return e.running
}
```

#### 2.5 代理管理器

**文件**: `internal/proxy/manager.go`
```go
package proxy

import (
    "fmt"
    "os"
    "path/filepath"
    "strconv"
)

type Manager struct {
    engine     *MihomoEngine
    configPath string
    pidFile    string
}

func NewManager(configDir string) *Manager {
    return &Manager{
        configPath: filepath.Join(configDir, "config.yaml"),
        pidFile:    filepath.Join(configDir, "clash-fish.pid"),
    }
}

func (m *Manager) Start() error {
    // 检查是否已运行
    if m.IsRunning() {
        return fmt.Errorf("service is already running")
    }

    // 创建引擎
    m.engine = NewMihomoEngine(m.configPath)

    // 启动引擎
    if err := m.engine.Start(); err != nil {
        return err
    }

    // 保存 PID
    return m.savePID()
}

func (m *Manager) Stop() error {
    if !m.IsRunning() {
        return fmt.Errorf("service is not running")
    }

    if err := m.engine.Stop(); err != nil {
        return err
    }

    // 删除 PID 文件
    return os.Remove(m.pidFile)
}

func (m *Manager) IsRunning() bool {
    _, err := os.Stat(m.pidFile)
    return err == nil
}

func (m *Manager) savePID() error {
    pid := os.Getpid()
    return os.WriteFile(m.pidFile, []byte(strconv.Itoa(pid)), 0644)
}
```

#### 2.6 权限检查

**文件**: `pkg/utils/permissions.go`
```go
package utils

import (
    "fmt"
    "os"
)

func CheckRoot() error {
    if os.Geteuid() != 0 {
        return fmt.Errorf("this command requires root privileges, please run with sudo")
    }
    return nil
}

func IsRoot() bool {
    return os.Geteuid() == 0
}
```

#### 2.7 实现 start 命令

**文件**: `cmd/clash-fish/start.go`
```go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/yourusername/clash-fish/internal/proxy"
    "github.com/yourusername/clash-fish/pkg/utils"
)

var startCmd = &cobra.Command{
    Use:   "start",
    Short: "Start clash-fish service",
    RunE:  runStart,
}

func runStart(cmd *cobra.Command, args []string) error {
    // 检查权限
    if err := utils.CheckRoot(); err != nil {
        return err
    }

    // 创建代理管理器
    manager := proxy.NewManager(getConfigDir())

    // 启动服务
    if err := manager.Start(); err != nil {
        return err
    }

    fmt.Println("✓ Clash-Fish started successfully")
    return nil
}

func init() {
    rootCmd.AddCommand(startCmd)
}
```

**检查点**:
- `sudo clash-fish start` 可以启动 mihomo
- PID 文件正确创建
- 日志显示启动信息

**Day 2 交付物**: 基础的 start/stop 功能可用

---

## Day 3: 系统集成与 VPN 检测

### 上午任务：完善启动/停止（3小时）

#### 3.1 信号处理

**文件**: `internal/service/daemon.go`
```go
package service

import (
    "os"
    "os/signal"
    "syscall"
    "github.com/rs/zerolog/log"
)

type Daemon struct {
    manager *proxy.Manager
}

func NewDaemon(manager *proxy.Manager) *Daemon {
    return &Daemon{manager: manager}
}

func (d *Daemon) Run() error {
    // 设置信号处理
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

    // 等待信号
    sig := <-sigCh
    log.Info().Msgf("Received signal: %v", sig)

    // 优雅关闭
    return d.manager.Stop()
}
```

#### 3.2 完善 stop 和 status 命令

**文件**: `cmd/clash-fish/stop.go`
```go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/yourusername/clash-fish/internal/proxy"
    "github.com/yourusername/clash-fish/pkg/utils"
)

var stopCmd = &cobra.Command{
    Use:   "stop",
    Short: "Stop clash-fish service",
    RunE:  runStop,
}

func runStop(cmd *cobra.Command, args []string) error {
    if err := utils.CheckRoot(); err != nil {
        return err
    }

    manager := proxy.NewManager(getConfigDir())

    if err := manager.Stop(); err != nil {
        return err
    }

    fmt.Println("✓ Clash-Fish stopped successfully")
    return nil
}

func init() {
    rootCmd.AddCommand(stopCmd)
}
```

**文件**: `cmd/clash-fish/status.go`
```go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/yourusername/clash-fish/internal/proxy"
)

var statusCmd = &cobra.Command{
    Use:   "status",
    Short: "Show clash-fish service status",
    RunE:  runStatus,
}

func runStatus(cmd *cobra.Command, args []string) error {
    manager := proxy.NewManager(getConfigDir())

    if manager.IsRunning() {
        fmt.Println("✓ Clash-Fish is running")
        // TODO: 显示更多状态信息
    } else {
        fmt.Println("✗ Clash-Fish is not running")
    }

    return nil
}

func init() {
    rootCmd.AddCommand(statusCmd)
}
```

---

### 下午任务：VPN 检测与系统集成（4小时）

#### 3.3 VPN 检测

**文件**: `internal/system/vpn.go`
```go
package system

import (
    "net"
    "strings"
    "github.com/rs/zerolog/log"
)

type VPNInfo struct {
    Active    bool
    Interface string
    IP        string
}

func DetectVPN() (*VPNInfo, error) {
    interfaces, err := net.Interfaces()
    if err != nil {
        return nil, err
    }

    for _, iface := range interfaces {
        // 检查 utun 接口
        if strings.HasPrefix(iface.Name, "utun") {
            addrs, err := iface.Addrs()
            if err != nil {
                continue
            }

            for _, addr := range addrs {
                ip, _, err := net.ParseCIDR(addr.String())
                if err != nil {
                    continue
                }

                // 检查是否是私网地址（VPN 特征）
                if isPrivateIP(ip) {
                    return &VPNInfo{
                        Active:    true,
                        Interface: iface.Name,
                        IP:        ip.String(),
                    }, nil
                }
            }
        }
    }

    return &VPNInfo{Active: false}, nil
}

func isPrivateIP(ip net.IP) bool {
    privateRanges := []string{
        "10.0.0.0/8",
        "172.16.0.0/12",
        "192.168.0.0/16",
    }

    for _, cidr := range privateRanges {
        _, ipnet, _ := net.ParseCIDR(cidr)
        if ipnet.Contains(ip) {
            return true
        }
    }

    return false
}
```

#### 3.4 路由表查看

**文件**: `internal/system/route.go`
```go
package system

import (
    "os/exec"
    "strings"
)

type RouteEntry struct {
    Destination string
    Gateway     string
    Interface   string
}

func GetRoutes() ([]RouteEntry, error) {
    cmd := exec.Command("netstat", "-nr")
    output, err := cmd.Output()
    if err != nil {
        return nil, err
    }

    var routes []RouteEntry
    lines := strings.Split(string(output), "\n")

    for _, line := range lines {
        // 解析路由表
        // TODO: 实现完整的路由解析逻辑
    }

    return routes, nil
}
```

#### 3.5 在启动时检测 VPN

更新 `internal/proxy/manager.go`:
```go
func (m *Manager) Start() error {
    // VPN 检测
    vpnInfo, err := system.DetectVPN()
    if err != nil {
        log.Warn().Err(err).Msg("Failed to detect VPN")
    } else if vpnInfo.Active {
        log.Info().
            Str("interface", vpnInfo.Interface).
            Str("ip", vpnInfo.IP).
            Msg("VPN detected, will coexist with proxy")
    }

    // ... 原有启动逻辑
}
```

**检查点**:
- VPN 检测功能正常
- 启动时显示 VPN 状态
- start/stop/status 命令完整可用

**Day 3 交付物**: 核心功能完整，VPN 共存机制验证通过

---

## Day 4: 订阅管理

### 上午任务：Profile 管理（3小时）

#### 4.1 Profile 数据结构

**文件**: `internal/subscription/profile.go`
```go
package subscription

import (
    "time"
)

type Profile struct {
    Name      string    `json:"name"`
    URL       string    `json:"url"`
    Path      string    `json:"path"`
    UpdatedAt time.Time `json:"updated_at"`
}

type ProfileManager struct {
    profilesDir string
    profiles    map[string]*Profile
    active      string
}

func NewProfileManager(profilesDir string) *ProfileManager {
    return &ProfileManager{
        profilesDir: profilesDir,
        profiles:    make(map[string]*Profile),
    }
}

func (pm *ProfileManager) Add(name, url string) error {
    // TODO: 实现添加 profile
    return nil
}

func (pm *ProfileManager) List() []*Profile {
    var list []*Profile
    for _, p := range pm.profiles {
        list = append(list, p)
    }
    return list
}

func (pm *ProfileManager) Switch(name string) error {
    // TODO: 实现切换 profile
    return nil
}

func (pm *ProfileManager) Delete(name string) error {
    // TODO: 实现删除 profile
    return nil
}
```

---

### 下午任务：订阅导入（4小时）

#### 4.2 HTTP 下载器

**文件**: `internal/subscription/downloader.go`
```go
package subscription

import (
    "io"
    "net/http"
    "time"
)

type Downloader struct {
    client *http.Client
}

func NewDownloader() *Downloader {
    return &Downloader{
        client: &http.Client{
            Timeout: 30 * time.Second,
        },
    }
}

func (d *Downloader) Download(url string) ([]byte, error) {
    resp, err := d.client.Get(url)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    return io.ReadAll(resp.Body)
}
```

#### 4.3 订阅导入器

**文件**: `internal/subscription/importer.go`
```go
package subscription

import (
    "gopkg.in/yaml.v3"
    "github.com/yourusername/clash-fish/internal/config"
)

type Importer struct {
    downloader *Downloader
}

func NewImporter() *Importer {
    return &Importer{
        downloader: NewDownloader(),
    }
}

func (i *Importer) Import(url string) (*config.Config, error) {
    // 下载订阅
    data, err := i.downloader.Download(url)
    if err != nil {
        return nil, err
    }

    // 解析为 Clash 配置
    var cfg config.Config
    if err := yaml.Unmarshal(data, &cfg); err != nil {
        return nil, err
    }

    return &cfg, nil
}
```

#### 4.4 Profile 命令实现

**文件**: `cmd/clash-fish/profile.go`
```go
package main

import (
    "fmt"
    "github.com/spf13/cobra"
    "github.com/yourusername/clash-fish/internal/subscription"
)

var profileCmd = &cobra.Command{
    Use:   "profile",
    Short: "Manage configuration profiles",
}

var profileAddCmd = &cobra.Command{
    Use:   "add <name> <url>",
    Short: "Add a new profile from subscription URL",
    Args:  cobra.ExactArgs(2),
    RunE:  runProfileAdd,
}

func runProfileAdd(cmd *cobra.Command, args []string) error {
    name := args[0]
    url := args[1]

    importer := subscription.NewImporter()

    // 导入订阅
    cfg, err := importer.Import(url)
    if err != nil {
        return err
    }

    // 保存配置
    // TODO: 保存到 profiles 目录

    fmt.Printf("✓ Profile '%s' added successfully\n", name)
    return nil
}

func init() {
    profileCmd.AddCommand(profileAddCmd)
    rootCmd.AddCommand(profileCmd)
}
```

**检查点**:
- `clash-fish profile add test https://xxx` 可用
- 配置文件正确下载和保存
- Profile 列表功能可用

**Day 4 交付物**: 订阅管理功能完整

---

## Day 5: 完善与 MVP 发布

### 上午任务：日志与监控（3小时）

#### 5.1 日志查看命令

**文件**: `cmd/clash-fish/log.go`
```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
    Use:   "log",
    Short: "Show logs",
}

var logTailCmd = &cobra.Command{
    Use:   "tail",
    Short: "Show real-time logs",
    RunE:  runLogTail,
}

func runLogTail(cmd *cobra.Command, args []string) error {
    logFile := filepath.Join(getConfigDir(), "logs", "clash-fish.log")

    // TODO: 实现实时日志查看（使用 tail -f 或类似机制）

    return nil
}

func init() {
    logCmd.AddCommand(logTailCmd)
    rootCmd.AddCommand(logCmd)
}
```

#### 5.2 增强 status 命令

更新 `cmd/clash-fish/status.go`:
```go
func runStatus(cmd *cobra.Command, args []string) error {
    manager := proxy.NewManager(getConfigDir())

    fmt.Println("=== Clash-Fish Status ===")

    // 服务状态
    if manager.IsRunning() {
        fmt.Println("Service:    ✓ Running")
    } else {
        fmt.Println("Service:    ✗ Not Running")
        return nil
    }

    // VPN 检测
    vpnInfo, _ := system.DetectVPN()
    if vpnInfo.Active {
        fmt.Printf("VPN:        ✓ Active (%s: %s)\n", vpnInfo.Interface, vpnInfo.IP)
    } else {
        fmt.Println("VPN:        ✗ Not Detected")
    }

    // 配置信息
    cfg, _ := config.Load()
    fmt.Printf("Mode:       %s\n", cfg.Mode)
    fmt.Printf("HTTP Port:  %d\n", cfg.Port)
    fmt.Printf("SOCKS Port: %d\n", cfg.SocksPort)

    return nil
}
```

---

### 下午任务：测试与文档（3小时）

#### 5.3 创建 Makefile

```makefile
.PHONY: build install uninstall clean test

BINARY_NAME=clash-fish
INSTALL_PATH=/usr/local/bin

build:
	go build -o $(BINARY_NAME) cmd/clash-fish/*.go

install: build
	sudo cp $(BINARY_NAME) $(INSTALL_PATH)/
	sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME)

uninstall:
	sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME)

clean:
	rm -f $(BINARY_NAME)
	go clean

test:
	go test -v ./...

run: build
	sudo ./$(BINARY_NAME) start
```

#### 5.4 编写 README

创建 `README.md`:
```markdown
# Clash-Fish

A lightweight transparent proxy tool for macOS based on mihomo core.

## Features

- ✅ Transparent proxy with TUN mode
- ✅ Automatic VPN coexistence
- ✅ Subscription import
- ✅ Multiple profiles support
- ✅ CLI interface

## Quick Start

### Installation

\`\`\`bash
make install
\`\`\`

### Usage

\`\`\`bash
# Initialize configuration
clash-fish config init

# Add subscription
clash-fish profile add myproxy https://example.com/clash

# Start service
sudo clash-fish start

# Check status
clash-fish status

# Stop service
sudo clash-fish stop
\`\`\`

## Requirements

- macOS 12+
- Root privileges for TUN mode

## License

MIT
```

#### 5.5 集成测试

创建测试清单 `TEST_CHECKLIST.md`:

```markdown
# MVP 测试清单

## 基础功能
- [ ] 编译成功
- [ ] config init 创建配置文件
- [ ] start 命令启动服务
- [ ] stop 命令停止服务
- [ ] status 显示正确状态

## 订阅功能
- [ ] profile add 导入订阅
- [ ] profile list 列出配置
- [ ] profile switch 切换配置

## VPN 共存
- [ ] 检测 VPN 连接
- [ ] VPN 流量正确路由
- [ ] 代理流量正确路由

## 错误处理
- [ ] 重复启动提示
- [ ] 无 root 权限提示
- [ ] 配置文件错误提示

## 文档
- [ ] README 完整
- [ ] 命令帮助信息完整
```

**Day 5 交付物**: MVP 1.0 版本完成

---

## Day 6-7: 测试、优化与文档

### Day 6: 全面测试

#### 功能测试
- [ ] 所有命令功能测试
- [ ] VPN 共存场景测试
- [ ] 订阅导入测试（多个订阅源）
- [ ] 错误场景测试
- [ ] 长时间运行稳定性测试

#### 性能测试
- [ ] 代理延迟测试
- [ ] 内存占用监控
- [ ] CPU 使用率监控

### Day 7: 优化与发布

#### 代码优化
- [ ] 错误处理完善
- [ ] 日志优化
- [ ] 代码重构

#### 文档完善
- [ ] 用户文档
- [ ] 故障排查文档
- [ ] 开发文档

#### 发布准备
- [ ] 版本标签
- [ ] Release Notes
- [ ] 安装脚本

---

## 检查清单总览

### MVP 必须功能（P0）
- [x] 项目初始化
- [ ] CLI 框架
- [ ] 配置管理
- [ ] Mihomo 集成
- [ ] start/stop/status 命令
- [ ] VPN 检测
- [ ] 订阅导入
- [ ] 基础文档

### 应该有功能（P1）
- [ ] Profile 切换
- [ ] 日志查看
- [ ] 权限检查
- [ ] 错误处理

### 可选功能（P2）
- [ ] 流量统计
- [ ] 代理测速
- [ ] 模式切换

---

## 开发环境设置

### 推荐工具
- **IDE**: VSCode / GoLand
- **调试**: Delve
- **格式化**: gofmt / goimports
- **Lint**: golangci-lint

### VSCode 配置

`.vscode/settings.json`:
```json
{
    "go.useLanguageServer": true,
    "go.lintTool": "golangci-lint",
    "go.lintOnSave": "package",
    "editor.formatOnSave": true
}
```

---

## 常见问题

### Q: 如何调试需要 root 权限的代码？
A: 使用 `sudo dlv debug` 或在 IDE 中配置 sudo 启动

### Q: Mihomo 依赖包太大怎么办？
A: 考虑使用二进制模式而不是库模式集成

### Q: 如何测试 VPN 共存？
A: 连接公司 VPN 后运行 clash-fish，检查路由表

---

**祝开发顺利！**
